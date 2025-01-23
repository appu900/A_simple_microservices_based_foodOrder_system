package handlers

import (
	"fmt"
	"log"
	"resturantService/database"
	"resturantService/model"
	"resturantService/utils"
	"sync"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var DishCollection *mongo.Collection
var once sync.Once

func IntiDishCollection() {
	once.Do(func() {
		DishCollection = database.GetCollection("dishes")
	})
}

func AddRestaurant(c *fiber.Ctx) error {
	var inputData struct {
		Name      string  `json:"name"`
		Address   string  `json:"address"`
		Photo     string  `json:"photo"`
		Lattitude float64 `json:"lattitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.BodyParser(&inputData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Inavlid Request Body",
		})
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "photo is required",
		})
	}

	if !utils.IsValidImageType(file.Header.Get("Content-Type")) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invaid image type",
		})
	}

	photoUrl, err := utils.UpLoadImageToS3(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Sommething went wrong in uploadig image",
		})
	}

	collection := database.GetCollection("restaurants")
	newRestaurantObject := model.NewRestaurant(inputData.Name, photoUrl, inputData.Address, inputData.Longitude, inputData.Lattitude)
	_, err = collection.InsertOne(c.Context(), newRestaurantObject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errorMessage": "something went wrommg in add resturant",
			"error":        err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok":   true,
		"data": newRestaurantObject,
	})
}

func AddDishes(c *fiber.Ctx) error {
	var inputData struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		Photo       string  `json:"photo"`
	}

	if err := c.BodyParser(&inputData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Request Body",
		})
	}

	if DishCollection == nil {
		IntiDishCollection()
	}

	restaurantId := c.Params("id")
	if restaurantId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "restaurant id is required",
		})
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "photo is required",
		})
	}

	if !utils.IsValidImageType(file.Header.Get("Content-Type")) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image type",
		})
	}

	photoUrl, err := utils.UpLoadImageToS3(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong in uploading image",
		})
	}

	newDishObject := model.NewDish(inputData.Name, inputData.Price, inputData.Description, photoUrl)
	_, err = DishCollection.InsertOne(c.Context(), newDishObject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong in adding dish",
		})
	}

	// update restaurant menu
	restaurantCollection := database.GetCollection("restaurants")
	primitiveRestaurantId, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid restaurant id",
		})
	}

	_, err = restaurantCollection.UpdateOne(c.Context(), bson.M{"_id": primitiveRestaurantId}, bson.M{"$push": bson.M{"menu": newDishObject.ID}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong in updating restaurant menu",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok":   true,
		"data": newDishObject,
	})

}

func GetAllMenu(c *fiber.Ctx) error {
	restaurantId := c.Params("id")
	if restaurantId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "restaurant id is required",
		})
	}

	restaurantCollection := database.GetCollection("restaurants")
	primitiveRestaurantId, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid restaurant id",
		})
	}

	// get all dishes from restaurant id

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": primitiveRestaurantId}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "dishes",
			"localField":   "menu",
			"foreignField": "_id",
			"as":           "menuDetails",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$menuDetails",
			"preserveNullAndEmptyArrays": true,
		}}},
	}

	cursor, err := restaurantCollection.Aggregate(c.Context(), pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	var results []bson.M
	if err := cursor.All(c.Context(), &results); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong in getting menu",
		})
	}

	fmt.Print(results)

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No menu found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":   true,
		"data": results,
	})
}

// Get all resturants in a given location with pegination
// the redius is in kilometers is 3 km only

func GetRestaurants(c *fiber.Ctx) error {
	var inputData struct {
		Longitude float64 `json:"longitude"`
		Lattitude float64 `json:"lattitude"`
	}
	if err := c.BodyParser(&inputData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var restaurantCollection = database.GetCollection("restaurants")
	query := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{inputData.Longitude, inputData.Lattitude},
				},
				"$maxDistance": 5000,
				 // 3 kilometers in meters
			},
		},
	}

	cursor, err := restaurantCollection.Find(c.Context(), query)
	if err != nil {
		log.Print("Error fetching restaurants By Location: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer cursor.Close(c.Context())

	var restaurants []model.Restaurant
	if err = cursor.All(c.Context(), &restaurants); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong while fetching restaurants",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":   true,
		"data": restaurants,
	})
}

package handlers

import (
	"resturantService/database"
	"resturantService/model"
	"resturantService/utils"

	"github.com/gofiber/fiber/v2"
)

func AddRestaurant(c *fiber.Ctx) error {
	var inputData struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Photo   string `json:"photo"`
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
	newRestaurantObject := model.NewRestaurant(inputData.Name, photoUrl, inputData.Address)
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

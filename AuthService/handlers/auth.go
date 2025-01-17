package handlers

import (
	"github.com/appu900/authservice/database"
	"github.com/appu900/authservice/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userCollection := database.GetCollection("users")
	_, err := userCollection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to register User"})
	}

	return c.JSON(fiber.Map{"message": "User register successfully"})
}

func Login(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Request body"})
	}

	userCollection := database.GetCollection("users")
	var foundUser models.User
	err := userCollection.FindOne(c.Context(), bson.M{"username": user.Username, "password": user.Password}).Decode(&foundUser)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Credentials"})
	}

	return nil

}



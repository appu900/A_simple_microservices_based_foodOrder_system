package handlers

import (
	"fmt"
	"time"
	"github.com/appu900/authservice/database"
	"github.com/appu900/authservice/models"
	"github.com/appu900/authservice/types"
	"github.com/appu900/authservice/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HandleUserRegistration(c *fiber.Ctx) error {
	userCollection := database.GetCollection("users")
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Request body",
		})
	}

	if input.Username == "" || input.Password == "" || input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields (username, password, email) are required",
		})
	}

	validationMessage := utils.Validate(input.Username, input.Password)
	if validationMessage != "Passed" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationMessage,
		})
	}

	existingUser, err := models.CheckIfUserExistsWithEmail(c, input.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "user already exists with this email",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	newUser := models.NewUser(input.Username, string(hashedPassword), input.Email)
	_, err = userCollection.InsertOne(c.Context(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create User",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "ok",
		"id":     newUser.ID.Hex(),
	})

}

func HandleLogin(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Request Body",
		})
	}

	if input.Email == "" || input.Password == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "all fields are required",
		})
	}

	userCollection := database.GetCollection("users")
	var user models.User
	err := userCollection.FindOne(c.Context(), bson.M{
		"email": input.Email,
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User Not found with this email",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Incorrect Password",
		})
	}

	now := time.Now()
	user.LastLogin = &now
	user.UpdatedAt = now

	_, err = userCollection.UpdateOne(
		c.Context(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"last_login": user.LastLogin,
			"updated_at": user.UpdatedAt,
		}},
	)

	if err != nil {
		fmt.Println("Failed to updated last login ")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretkey := []byte("hello_brother_key")

	tokenString, err := token.SignedString(secretkey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong with token creation",
		})
	}

	return c.JSON(types.LoginResponse{
		Token:     tokenString,
		TokenType: "Bearer",
	})
}

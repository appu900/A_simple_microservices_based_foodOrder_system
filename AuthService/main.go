package main

import (
	"log"

	"github.com/appu900/authservice/database"
	"github.com/appu900/authservice/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	if err := database.Connect(); err != nil {
		log.Fatal("Something went wrong while connecting to the database", err)
	}

	app.Post("/api/register", handlers.HandleUserRegistration)
	app.Post("/api/login", handlers.HandleLogin)

	log.Fatal(app.Listen(":3000"))
}

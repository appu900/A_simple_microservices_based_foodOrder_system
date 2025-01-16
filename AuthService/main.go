package main

import (
	"log"
	"github.com/appu900/authservice/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	if err := database.Connect(); err != nil {
		log.Fatal("Something went wrong while connecting to the database", err)
	}

	app.Get("/pingme", func(c *fiber.Ctx) error {
		return c.SendString("Pinged auth service")
	})

	log.Fatal(app.Listen(":3000"))
}

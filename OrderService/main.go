package main

import (
	"log"
	"github.com/appu900/OrderService/database"
	"github.com/appu900/OrderService/handler"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	// connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Error connecting to DB : %v", err)
	}

	// routes for testing purpose
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! form order service")
	})

	app.Get("/api/health", handler.HealthCheck)
	log.Fatal(app.Listen(":3002"))
}

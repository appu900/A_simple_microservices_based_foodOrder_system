package main

import (
	"log"
	"resturantService/database"
	"resturantService/handlers"
	"resturantService/model"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// connect to database
	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Error connecting to DB")
	}

	// make index for resturant service
	model.CreateGeospatialIndex()

	app := fiber.New()

	app.Post("/api/restaurant", handlers.AddRestaurant)
	app.Post("/api/restaurant/:id/dish", handlers.AddDishes)
	app.Get("/api/restaurant/:id/menu", handlers.GetAllMenu)
	app.Get("/api/restaurant/userLocation", handlers.GetRestaurants)

	log.Fatal(app.Listen(":3001"))
}

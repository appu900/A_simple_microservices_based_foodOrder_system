package main

import (
	"log"
	"github.com/appu900/authservice/database"
	"github.com/appu900/authservice/handlers"
	"github.com/appu900/authservice/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	if err := database.Connect(); err != nil {
		log.Fatal("Something went wrong while connecting to the database", err)
	}

	app.Post("/api/register", handlers.HandleUserRegistration)
	app.Post("/api/login", handlers.HandleLogin)
	app.Post("/api/task/create", middleware.AuthMiddleware(), handlers.HandleCreateTask)
	app.Get("/api/task/get", middleware.AuthMiddleware(), handlers.HandleGetAlltasksOfUser)
	app.Put("/api/task/update/:id", middleware.AuthMiddleware(), handlers.HandleUpdateTask)

	log.Fatal(app.Listen(":3000"))
}

package handler

import "github.com/gofiber/fiber/v2"

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Healthy! Order Service is up and running. :)")
}


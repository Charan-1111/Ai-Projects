package handlers

import "github.com/gofiber/fiber/v3"

func sendSuccess(c fiber.Ctx, status int, message string, payload fiber.Map) error {
	response := fiber.Map{
		"code":    0,
		"message": message,
	}
	for key, value := range payload {
		response[key] = value
	}
	return c.Status(status).JSON(response)
}

func sendError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"code":    1,
		"message": message,
	})
}

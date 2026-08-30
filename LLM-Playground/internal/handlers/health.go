package handlers

import "github.com/gofiber/fiber/v3"

func (h *Handlers) HealthCheck(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Server is health",
	})
}

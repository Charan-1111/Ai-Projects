package handlers

import "github.com/gofiber/fiber/v3"

func (h *Handlers) HealthCheck(c fiber.Ctx) error {
	return sendSuccess(c, fiber.StatusOK, "Server is healthy", fiber.Map{})
}

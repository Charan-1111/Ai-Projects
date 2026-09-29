package handler

import (
	"tool-calling/internal/models"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) GetTime(c fiber.Ctx) error {
	var req models.GetTimeRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code" : 1, "message" : "Invalid request body"})
	}

	
	return nil
}

package handler

import (
	"tool-calling/internal/models"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) GetCurrentTime(c fiber.Ctx) error {
	var currentTimeArgs models.TimeInput
	if err := c.Bind().Body(&currentTimeArgs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid request body"})
	}

	currentTime, err := h.services.GetCurrentTime(c.Context(), currentTimeArgs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 2, "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "time fetched successful", "data": currentTime})
}

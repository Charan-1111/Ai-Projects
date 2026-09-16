package handlers

import (
	"semantic-search/internal/models"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) CompareSimilarity(c fiber.Ctx) error {
	var reqBody models.Compare

	if err := c.Bind().Body(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid request body"})
	}

	similarity, err := h.service.CompareSimilarity(c.Context(), reqBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code" : 1, "message" : err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code" : 0, "message" : "similarity is done", "similarity" : similarity})
}

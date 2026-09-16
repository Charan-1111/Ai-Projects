package handlers

import (
	"semantic-search/internal/models"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) EmbedText(c fiber.Ctx) error {
	var req models.EmbedRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid request body"})
	}

	embedding, err := h.service.EmbedText(c.Context(), req.Text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": 1, "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":     0,
		"message":  "Text embedded successfully",
		"contents": models.EmbedResponse{Text: req.Text, Embedding: embedding},
	})
}

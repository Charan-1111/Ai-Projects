package handlers

import (
	"semantic-search/internal/models"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) InjectDocument(c fiber.Ctx) error {
	var docReq models.Document

	if err := c.Bind().Body(&docReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid request body"})
	}

	docResponse, err := h.service.EmbedDocument(c.Context(), docReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": 1, "message": "Error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code" : 0, "message" : "Embedding successful", "contents" : docResponse})
}

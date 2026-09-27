package handlers

import (
	"ask-my-documents/internal/models"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) AskDocuments(c fiber.Ctx) error {
	var req models.AskDocument
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid request body"})
	}

	semanticResponse, err := h.service.AskDocuments(c.Context(), req)
	if err != nil {

	}
	return c.Status(fiber.StatusOK).JSON(semanticResponse)
}

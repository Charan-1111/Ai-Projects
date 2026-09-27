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

	if req.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "prompt is required"})
	}
	if req.NoOfDocs < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "noOfDocs must be >= 0"})
	}
	if req.MinimumScore < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "minimumScore must be >= 0"})
	}

	llmResponse, err := h.service.AskDocuments(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": 1, "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(llmResponse)
}

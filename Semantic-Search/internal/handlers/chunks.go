package handlers

import (
	"semantic-search/internal/models"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) UploadDocuments(c fiber.Ctx) error {
	var req []models.Document
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid request body"})
	}

	if err := h.service.MakeChunksAndUpload(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": 1, "message": "Error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Chunks uploaded successfully"})
}

func (h *Handlers) SearchChunkedDocuments(c fiber.Ctx) error {
	return nil
}

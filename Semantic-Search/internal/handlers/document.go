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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Embedding successful", "contents": docResponse})
}

func (h *Handlers) UpdateDocument(c fiber.Ctx) error {
	var docReq models.Document

	if err := c.Bind().Body(&docReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid request body"})
	}

	docResponse, updated, err := h.service.UpdateDocument(c.Context(), c.Params("docId"), docReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": 1, "message": "Error"})
	}
	if !updated {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": 1, "message": "Document not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Document updated successfully", "contents": docResponse})
}

func (h *Handlers) DeleteDocument(c fiber.Ctx) error {
	deleted, err := h.service.DeleteDocument(c.Context(), c.Params("docId"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": 1, "message": "Error"})
	}
	if !deleted {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": 1, "message": "Document not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Document deleted successfully"})
}

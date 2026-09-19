package handlers

import (
	"semantic-search/internal/models"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) SearchDocuments(c fiber.Ctx) error {
	// TODO: Need to provide proper error messages and status codes
	var req models.SearchRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code" : 1, "message" : "Invalid request body"})
	}

	searchResponse, err := h.service.SearchDocuments(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code" : 1, "message" : "Something went wrong"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code" : 0, "message" : "Documents found", "content" : searchResponse})
}

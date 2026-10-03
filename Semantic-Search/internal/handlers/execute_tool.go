package handlers

import (
	"semantic-search/internal/tools"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) ExecuteTools(c fiber.Ctx) error {
	var toolReq tools.ExecuteToolRequest
	if err := c.Bind().Body(&toolReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid request body"})
	}

	result, err := h.service.ExecuteTools(c.Context(), toolReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": 1, "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

package handlers

import "github.com/gofiber/fiber/v3"

func (h *Handlers) ListTools(c fiber.Ctx) error {
	tools := h.service.ListTools(c.Context())

	return c.Status(fiber.StatusOK).JSON(tools)
}

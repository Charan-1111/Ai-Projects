package handler

import "github.com/gofiber/fiber/v3"

func (h *Handler) ListTools(c fiber.Ctx) error {
	tools := h.services.ListTools()
	return c.Status(fiber.StatusOK).JSON(tools)
}

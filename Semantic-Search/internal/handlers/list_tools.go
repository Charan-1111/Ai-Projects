package handlers

import "github.com/gofiber/fiber/v3"

func (h *Handlers) ListTools(c fiber.Ctx) error {
	h.log.Log.Info().Msg("Semantic Tool calling")
	tools := h.service.ListTools(c.Context())

	return c.Status(fiber.StatusOK).JSON(tools)
}

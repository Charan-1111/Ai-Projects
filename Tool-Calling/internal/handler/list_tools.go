package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) ListTools(c fiber.Ctx) error {
	fmt.Println("Tool calling")
	tools := h.services.ListTools()
	return c.Status(fiber.StatusOK).JSON(tools)
}

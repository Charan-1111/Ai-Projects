package handlers

import "github.com/gofiber/fiber/v3"

func (h *Handlers) AvailableModels(c fiber.Ctx) error {
	models, err := h.Services.AvailableModels()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code" : 1, "message" : err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code" : 0, "message" : "Model retrieval successful", "available_models" : models})
}

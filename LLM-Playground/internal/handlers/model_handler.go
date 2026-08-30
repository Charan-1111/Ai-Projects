package handlers

import "github.com/gofiber/fiber/v3"

func (h *Handlers) AvailableModels(c fiber.Ctx) error {
	models, err := h.Services.AvailableModels()
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return sendSuccess(c, fiber.StatusOK, "Model retrieval successful", fiber.Map{"available_models": models})
}

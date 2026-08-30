package handlers

import (
	"llm-playground/internal/models"
	"llm-playground/internal/validation"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) GenerateResponse(c fiber.Ctx) error {
	var request models.PromptRequest

	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid Request Body"})
	}

	if err := validation.ValidatePromptRequest(&request, h.config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": err.Error()})
	}

	requestId := c.Locals("requestId").(string)

	modelResponse, err := h.Services.ResponseGeneration(c.Context(), requestId, &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": 1, "message": "Failed to generate the response : " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Response generated successfully", "response": modelResponse})
}

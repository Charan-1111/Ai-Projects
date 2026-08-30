package handlers

import (
	"llm-playground/internal/models"
	"llm-playground/internal/validation"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) GenerateResponse(c fiber.Ctx) error {
	var request models.PromptRequest

	if err := c.Bind().Body(&request); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid Request Body")
	}

	if err := validation.ValidatePromptRequest(&request, h.config); err != nil {
		return sendError(c, fiber.StatusBadRequest, err.Error())
	}

	requestId := c.Locals("requestId").(string)

	modelResponse, err := h.Services.ResponseGeneration(c.Context(), requestId, &request)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to generate the response: "+err.Error())
	}

	return sendSuccess(c, fiber.StatusOK, "Response generated successfully", fiber.Map{"response": modelResponse})
}

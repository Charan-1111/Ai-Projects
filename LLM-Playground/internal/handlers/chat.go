package handlers

import (
	"llm-playground/internal/models"
	"llm-playground/internal/validation"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) Chat(c fiber.Ctx) error {
	var request models.PromptRequest
	if err := c.Bind().Body(&request); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid Request Body")
	}

	if err := validation.ValidatePromptRequest(&request, h.config); err != nil {
		return sendError(c, fiber.StatusBadRequest, err.Error())
	}

	requestId := c.Locals("requestId").(string)
	sessionId := c.Locals("sessionId").(string)


	
	return nil
}

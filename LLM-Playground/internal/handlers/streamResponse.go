package handlers

import (
	"llm-playground/internal/models"
	"llm-playground/internal/validation"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/sse"
)

func (h *Handlers) GenerateStreamResponse(c fiber.Ctx, stream *sse.Stream) error {
	var request models.PromptRequest

	if err := c.Bind().Body(&request); err != nil {
		return stream.Event(sse.Event{Name: "error", Data: fiber.Map{"code": 1, "message": "Invalid request body"}})
	}

	if err := validation.ValidatePromptRequest(&request, h.config); err != nil {
		return stream.Event(sse.Event{Name: "error", Data: fiber.Map{"code": 1, "message": err.Error()}})
	}

	requestId := c.Locals("requestId").(string)

	chunks, errs, err := h.Services.ResponseStreamGeneration(stream.Context(), &request)
	if err != nil {
		return stream.Event(sse.Event{Name: "error", Data: fiber.Map{"code": 1, "message": err.Error()}})
	}

	for chunks != nil || errs != nil {
		select {
		case chunk, ok := <-chunks:
			if !ok {
				chunks = nil
				continue
			}

			if evErr := stream.Event(sse.Event{
				Name: "chunk",
				Data: fiber.Map{
					"request_id":    requestId,
					"delta":         chunk.Delta,
					"input_tokens":  chunk.InputTokens,
					"output_tokens": chunk.OutputTokens,
					"finish_reason": chunk.FinishReason,
				},
			}); evErr != nil {
				return evErr
			}

		case streamErr, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}

			statusCode := providerErrorStatus(streamErr)
			return stream.Event(sse.Event{Name: "error", Data: fiber.Map{"code": statusCode, "message": streamErr.Error()}})
		}
	}

	return nil
}

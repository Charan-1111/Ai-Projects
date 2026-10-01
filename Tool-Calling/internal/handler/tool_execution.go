package handler

import (
	"tool-calling/internal/tools"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) ExecuteTool(c fiber.Ctx) error {
	var req tools.ExecuteToolRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(tools.ToolExecutionResponse{
			Error: &tools.ToolError{
				Code:    1,
				Message: "Invalid request body",
			},
		})
	}

	toolResult, err := h.services.ToolExecution(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(tools.ToolExecutionResponse{
			Error: &tools.ToolError{
				Code:    1,
				Message: "Failed to execute tool",
			},
		})
	}
	return c.Status(fiber.StatusOK).JSON(tools.ToolExecutionResponse{
		Result: toolResult.Result,
	})
}

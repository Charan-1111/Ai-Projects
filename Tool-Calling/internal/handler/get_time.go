package handler

import (
	"context"
	"fmt"
	"os"
	"strings"
	"tool-calling/internal/models"
	"tool-calling/internal/tools"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/genai"
)

func (h *Handler) GetTime(c fiber.Ctx) error {
	var req models.GetTimeRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 1, "message": "Invalid request body"})
	}

	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 2, "message": "Message is required"})
	}

	response, err := handleTimeQuestion(req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": 3, "message": err.Error()})
	}

	return c.JSON(response)
}

func handleTimeQuestion(message string) (models.GetTimeResponse, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	if apiKey == "" {
		return models.GetTimeResponse{}, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return models.GetTimeResponse{}, fmt.Errorf("failed to initialize Gemini client: %w", err)
	}

	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				tools.CurrentTimeDeclaration(),
			},
		}},
	}

	response, err := client.Models.GenerateContent(
		ctx,
		"gemini-3.5-flash-lite",
		[]*genai.Content{genai.NewContentFromText(message, genai.RoleUser)},
		config,
	)
	if err != nil {
		return models.GetTimeResponse{}, fmt.Errorf("gemini request failed: %w", err)
	}

	for _, candidate := range response.Candidates {
		if candidate == nil || candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			if part == nil || part.FunctionCall == nil {
				continue
			}
			if part.FunctionCall.Name != "get_current_time" {
				continue
			}

			timezone, _ := part.FunctionCall.Args["timezone"].(string)
			if timezone == "" {
				return models.GetTimeResponse{
					Answer:        "Please provide a timezone so I can fetch the current time.",
					ToolsExecutes: []string{"get_current_time"},
				}, nil
			}

			result, toolErr := tools.GetCurrentTime(tools.CurrentTimeArgs{Timezone: timezone})
			if toolErr != nil {
				return models.GetTimeResponse{}, fmt.Errorf("timezone tool failed: %w", toolErr)
			}

			return models.GetTimeResponse{
				Answer:        fmt.Sprintf("The current time in %s is %s.", result.Timezone, result.CurrentTime),
				ToolsExecutes: []string{"get_current_time"},
			}, nil
		}
	}

	if response != nil && response.Text() != "" {
		return models.GetTimeResponse{Answer: response.Text()}, nil
	}

	return models.GetTimeResponse{Answer: "I couldn't find an answer for that question."}, nil
}

package services

import (
	"context"
	"fmt"
	"llm-playground/internal/models"
	"llm-playground/internal/pricing"
	"time"

	"google.golang.org/genai"
)

func (s *Services) ResponseGeneration(ctx context.Context, requestId string, request *models.PromptRequest) (models.ModelResponse, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		30*time.Second,
	)
	defer cancel()

	start := time.Now()
	response, err := s.client.Models.GenerateContent(
		ctx,
		request.Model,
		genai.Text(request.Prompt),
		nil,
	)

	if err != nil {
		return models.ModelResponse{}, fmt.Errorf("Failed to generate response : %w", err)
	}

	latencyMs := time.Since(start).Milliseconds()

	usage, totalCost := pricing.PriceCalculator(s.config.AvailableModels[request.ModelId], response)

	modelResponse := models.ModelResponse{
		RequestId: requestId,
		Model:     request.Model,
		Response: response.Text(),
		Usage: usage,
		EstimatedCostUsed: totalCost,
		FinishReason: "stop",
		LatencyMs: latencyMs,
	}

	return modelResponse, nil
}

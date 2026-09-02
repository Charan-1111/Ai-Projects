package services

import (
	"context"
	"fmt"
	"llm-playground/internal/models"
	"llm-playground/internal/provider"
)

func (s *Services) ResponseStreamGeneration(ctx context.Context, request *models.PromptRequest) (<-chan provider.StreamChunk, <-chan error, error) {
	if request == nil {
		return nil, nil, fmt.Errorf("request body is required")
	}

	modelConfig, ok := s.resolveModelConfig(request)
	if !ok {
		return nil, nil, fmt.Errorf("model configuration not found for request model %q and model_id %q", request.Model, request.ModelId)
	}

	input, _ := provider.BuildGenerateInput(modelConfig, request)

	chunks, errs := s.provider.GenerateStream(ctx, input)
	return chunks, errs, nil
}

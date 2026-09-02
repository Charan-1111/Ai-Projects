package services

import (
	"context"
	"fmt"
	"llm-playground/internal/models"
	"llm-playground/internal/pricing"
	"llm-playground/internal/provider"
	"time"
)

func (s *Services) ResponseGeneration(ctx context.Context, requestId string, request *models.PromptRequest) (models.ModelResponse, error) {
	if request == nil {
		return models.ModelResponse{}, fmt.Errorf("request body is required")
	}

	modelConfig, ok := s.resolveModelConfig(request)
	if !ok {
		return models.ModelResponse{}, fmt.Errorf("model configuration not found for request model %q and model_id %q", request.Model, request.ModelId)
	}

	input, _ := provider.BuildGenerateInput(modelConfig, request)

	start := time.Now()
	response, err := s.provider.Generate(ctx, input)
	if err != nil {
		return models.ModelResponse{}, err
	}
	latencyMs := time.Since(start).Milliseconds()

	usage, totalCost := pricing.PriceCalculator(modelConfig, response)

	modelResponse := models.ModelResponse{
		RequestId:         requestId,
		Model:             modelConfig.ProviderModel,
		Response:          response.Text,
		Usage:             usage,
		EstimatedCostUsed: totalCost,
		FinishReason:      response.FinishReason,
		LatencyMs:         latencyMs,
	}

	return modelResponse, nil
}

func (s *Services) resolveModelConfig(request *models.PromptRequest) (models.ModelConfig, bool) {
	if request == nil || s == nil || s.config == nil {
		return models.ModelConfig{}, false
	}

	if request.ModelId != "" {
		if modelConfig, ok := s.config.AvailableModels[request.ModelId]; ok {
			if request.Model == "" {
				request.Model = modelConfig.ProviderModel
			}
			return modelConfig, true
		}
	}

	if request.Model != "" {
		for id, modelConfig := range s.config.AvailableModels {
			if modelConfig.ProviderModel == request.Model || id == request.Model {
				if request.ModelId == "" {
					request.ModelId = id
				}
				request.Model = modelConfig.ProviderModel
				return modelConfig, true
			}
		}
	}

	if s.config.DefaultModel != "" {
		for id, modelConfig := range s.config.AvailableModels {
			if modelConfig.ProviderModel == s.config.DefaultModel || id == s.config.DefaultModel {
				if request.ModelId == "" {
					request.ModelId = id
				}
				if request.Model == "" {
					request.Model = modelConfig.ProviderModel
				}
				return modelConfig, true
			}
		}
	}

	return models.ModelConfig{}, false
}

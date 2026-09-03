package services

import (
	"context"
	"fmt"
	"llm-playground/internal/models"
	"llm-playground/internal/pricing"
	"llm-playground/internal/provider"
	"llm-playground/internal/utils"
	"math/rand/v2"
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
	if s.config.Retries.MaxAttempts <= 0 || s.config.Retries.InitialDelayMs < 0 || s.config.Retries.MaximumDelayMs < 0 || s.config.Retries.MaximumDelayMs < s.config.Retries.InitialDelayMs {
		return models.ModelResponse{}, fmt.Errorf("invalid retry configuration")
	}

	input, _ := provider.BuildGenerateInput(modelConfig, request)

	var (
		response   *provider.GenerateResponse
		statusCode int
		err        error
	)

	start := time.Now()
	delay := time.Duration(s.config.Retries.InitialDelayMs) * time.Millisecond

	for attempt := 1; attempt <= s.config.Retries.MaxAttempts; attempt++ {
		response, statusCode, err = s.provider.Generate(ctx, input)
		if err == nil || !utils.IsRetryableError(statusCode) {
			break
		}

		if attempt == s.config.Retries.MaxAttempts {
			return models.ModelResponse{}, fmt.Errorf("failed to generate response after %d attempts: %w", attempt, err)
		}

		fmt.Printf("Attempt %d failed: %v. Retrying after up to %d ms...\n", attempt, err, delay.Milliseconds())

		withDuration := time.Duration(
			rand.Int64N(int64(delay) + 1),
		)

		timer := time.NewTimer(withDuration)

		select {
		case <-timer.C:
			nextDelay := delay * 2
			maximumDelay := time.Duration(s.config.Retries.MaximumDelayMs) * time.Millisecond
			if nextDelay > maximumDelay {
				nextDelay = maximumDelay
			}
			delay = nextDelay
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return models.ModelResponse{}, fmt.Errorf("context canceled during retry: %w", ctx.Err())
		}
	}
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

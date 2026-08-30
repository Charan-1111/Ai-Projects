package validation

import (
	"fmt"
	"llm-playground/internal/config"
	"llm-playground/internal/models"
	"strings"
)

func ValidatePromptRequest(request *models.PromptRequest, cfg *config.Configuration) error {
	if request == nil {
		return fmt.Errorf("request body is required")
	}

	if strings.TrimSpace(request.Prompt) == "" {
		return fmt.Errorf("prompt is required")
	}

	if request.Temperature < 0 || request.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	if request.MaxOutputTokens < 1 {
		return fmt.Errorf("max_output_tokens must be greater than 0")
	}

	if cfg == nil {
		return fmt.Errorf("server configuration is not available")
	}

	if request.ModelId != "" {
		modelConfig, ok := cfg.AvailableModels[request.ModelId]
		if !ok {
			return fmt.Errorf("invalid model_id: %s", request.ModelId)
		}
		if request.Model == "" {
			request.Model = modelConfig.ProviderModel
		}
		if request.MaxOutputTokens > int64(modelConfig.MaxOutputTokens) {
			return fmt.Errorf("max_output_tokens exceeds the configured limit for model %s", request.ModelId)
		}
		return nil
	}

	if request.Model != "" {
		for modelID, modelConfig := range cfg.AvailableModels {
			if modelConfig.ProviderModel == request.Model || modelID == request.Model {
				if request.ModelId == "" {
					request.ModelId = modelID
				}
				if request.MaxOutputTokens > int64(modelConfig.MaxOutputTokens) {
					return fmt.Errorf("max_output_tokens exceeds the configured limit for model %s", modelID)
				}
				return nil
			}
		}
		return fmt.Errorf("invalid model: %s", request.Model)
	}

	if cfg.DefaultModel != "" {
		request.Model = cfg.DefaultModel
		for modelID, modelConfig := range cfg.AvailableModels {
			if modelConfig.ProviderModel == cfg.DefaultModel || modelID == cfg.DefaultModel {
				request.ModelId = modelID
				if request.MaxOutputTokens > int64(modelConfig.MaxOutputTokens) {
					return fmt.Errorf("max_output_tokens exceeds the configured limit for model %s", modelID)
				}
				return nil
			}
		}
	}

	return fmt.Errorf("model or model_id is required")
}

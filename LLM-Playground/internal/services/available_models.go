package services

import (
	"fmt"
	"llm-playground/internal/models"
)

func (s *Services) AvailableModels() (*models.AvailableModelsResponse, error) {
	if len(s.config.AvailableModels) == 0 {
		return nil, fmt.Errorf("No models are available currently. Please try later")
	}

	availableModels := make([]models.AvailableModels, 0)

	for id, model := range s.config.AvailableModels {
		availableModel := models.AvailableModels{}
		availableModel.Id = id
		availableModel.DisplayName = model.DisplayName
		availableModel.MaxOutputTokens = model.MaxOutputTokens
		availableModel.SupportStreaming = model.SupportsStreaming

		availableModels = append(availableModels, availableModel)
	}

	return &models.AvailableModelsResponse{
		DefaultModel:    s.config.DefaultModel,
		AvailableModels: availableModels,
	}, nil
}

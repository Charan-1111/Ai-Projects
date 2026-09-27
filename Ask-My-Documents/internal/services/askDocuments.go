package services

import (
	"ask-my-documents/internal/models"
	"context"

	"github.com/bytedance/sonic"
)

func (s *Service) AskDocuments(ctx context.Context, req models.AskDocument) (models.SemanticResponse, error) {
	var semanticResponse models.SemanticResponse
	apiClient := s.apiFactory.Create(
		s.config.ExternalApis.SemanticSearch,
		"POST",
		map[string]string{},
		map[string]any{
			"query":        req.Prompt,
			"noOfDocs":     5,
			"minimumScore": 0.7,
			"filters": map[string]string{
				"category":   "golang",
				"difficulty": "beginner",
			},
		},
		map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		s.clients,
	)

	respDataBytes, err := apiClient.ApiCall(ctx)
	if err != nil {

	}

	err = sonic.Unmarshal(respDataBytes, &semanticResponse)
	if err != nil {

	}

	return semanticResponse, nil
}

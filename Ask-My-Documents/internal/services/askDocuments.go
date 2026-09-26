package services

import (
	"ask-my-documents/internal/models"
	"context"
)

func (s *Service) AskDocuments(ctx context.Context, req models.AskDocument) {
	// get the related chunks for the given prompt
	apiClient := s.apiFactory.Create(
		s.config.ExternalApis.SemanticSearch,
		"POST",
		map[string]string{},
		map[string]any{
			"query":    req.Prompt,
			"noOfDocs": 5,
		},
		map[string]string{},
		s.clients,
	)

	apiClient.ApiCall(ctx)
}

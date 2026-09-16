package services

import (
	"context"
	"semantic-search/internal/models"
)

func (s *Service) CompareSimilarity(ctx context.Context, reqBody models.Compare) (float64, error) {
	similarity, err := s.llmProvider.CompareSimilarity(ctx, reqBody.Original, reqBody.Compare)
	if err != nil {
		return 0, err
	}

	return similarity, nil
}

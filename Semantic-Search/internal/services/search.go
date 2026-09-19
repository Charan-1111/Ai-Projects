package services

import (
	"context"
	"semantic-search/internal/models"
)

const defaultSearchDocumentLimit = 10

func (s *Service) SearchDocuments(ctx context.Context, req models.SearchRequest) (models.SearchResponse, error) {
	if req.NoofDocs <= 0 {
		req.NoofDocs = defaultSearchDocumentLimit
	}

	queryEmbed, err := s.llmProvider.EmbedText(ctx, req.Query)
	if err != nil {
		return models.SearchResponse{}, err
	}

	similarDocuments, err := s.dbStore.SearchSimilarDocuments(ctx, queryEmbed, req.NoofDocs)
	if err != nil {
		return models.SearchResponse{}, err
	}

	searchResponse := models.SearchResponse{
		Query:            req.Query,
		SimilarDocuments: similarDocuments,
		ResultCount:      len(similarDocuments),
	}

	return searchResponse, nil
}

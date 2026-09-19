package services

import (
	"context"
	"semantic-search/internal/models"

	"github.com/google/uuid"
)

func (s *Service) EmbedText(ctx context.Context, text string) ([]float32, error) {
	return s.llmProvider.EmbedText(ctx, text)
}

type Document struct {
	Document     models.Document `json:"document"`
	DocEmbedding []float32       `json:"docEmbedding"`
}

func (s *Service) EmbedDocument(ctx context.Context, reqBody models.Document) (models.DocumentResponse, error) {
	embeddingFields := reqBody.DocTitle + " : " + reqBody.DocContent

	embed, err := s.llmProvider.EmbedText(ctx, embeddingFields)
	if err != nil {
		return models.DocumentResponse{}, err
	}

	docID := uuid.NewString()
	reqBody.DocId = docID

	err = s.dbStore.SaveDocument(ctx, reqBody, embed)
	if err != nil {
		return models.DocumentResponse{}, err
	}

	docResposne := models.DocumentResponse{
		DocId:           docID,
		DocTitle:        reqBody.DocTitle,
		EmbeddingStatus: "completed",
	}

	return docResposne, nil
}

func (s *Service) EmbedMultiDocument(ctx context.Context, req []models.Document) ([]models.DocumentResponse, error) {
	responses := make([]models.DocumentResponse, 0, len(req))

	for _, doc := range req {
		response, err := s.EmbedDocument(ctx, doc)
		if err != nil {
			return nil, err
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (s *Service) UpdateDocument(ctx context.Context, docID string, reqBody models.Document) (models.DocumentResponse, bool, error) {
	embeddingFields := reqBody.DocTitle + " : " + reqBody.DocContent

	embed, err := s.llmProvider.EmbedText(ctx, embeddingFields)
	if err != nil {
		return models.DocumentResponse{}, false, err
	}

	reqBody.DocId = docID
	updated, err := s.dbStore.UpdateDocument(ctx, reqBody, embed)
	if err != nil {
		return models.DocumentResponse{}, false, err
	}
	if !updated {
		return models.DocumentResponse{}, false, nil
	}

	return models.DocumentResponse{
		DocId:           docID,
		DocTitle:        reqBody.DocTitle,
		EmbeddingStatus: "completed",
	}, true, nil
}

func (s *Service) DeleteDocument(ctx context.Context, docID string) (bool, error) {
	return s.dbStore.DeleteDocument(ctx, docID)
}

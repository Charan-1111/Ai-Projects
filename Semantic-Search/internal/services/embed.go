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

	*s.documents = append(*s.documents, Document{
		Document:     reqBody,
		DocEmbedding: embed,
	})

	docResposne := models.DocumentResponse{
		DocId:           docID,
		DocTitle:        reqBody.DocTitle,
		EmbeddingStatus: "completed",
	}

	return docResposne, nil
}

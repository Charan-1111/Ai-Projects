package services

import (
	"context"
	"fmt"
	"semantic-search/internal/models"
)

func (s *Service) MakeChunksAndUpload(ctx context.Context, req []models.Document) error {
	for _, doc := range req {
		docResponse, err := s.EmbedDocument(ctx, doc)
		if err != nil {
			return fmt.Errorf("embed document: %w", err)
		}

		chunks, err := s.wordChunker.Split(ctx, doc.DocContent)
		if err != nil {
			return fmt.Errorf("split document into chunks: %w", err)
		}

		for index := range chunks {
			chunkEmbed, err := s.llmProvider.EmbedText(ctx, chunks[index].Content)
			if err != nil {
				return fmt.Errorf("embed chunk %d: %w", index, err)
			}

			chunks[index].ChunkEmbed = chunkEmbed
			chunks[index].DocId = docResponse.DocId
		}

		if err := s.dbStore.UploadChunks(ctx, chunks); err != nil {
			return fmt.Errorf("upload chunks: %w", err)
		}
	}

	return nil
}

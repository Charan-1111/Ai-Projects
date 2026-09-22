package database

import (
	"context"
	"fmt"
	"semantic-search/internal/chunks"
	"semantic-search/internal/models"

	"github.com/pgvector/pgvector-go"
)

type Repository interface {
	Create(ctx context.Context) error
	SaveDocument(ctx context.Context, document models.Document, embedding []float32) error
	UpdateDocument(ctx context.Context, document models.Document, embedding []float32) (bool, error)
	UpdateDocumentWithChunks(ctx context.Context, document models.Document, embedding []float32, documentChunks []chunks.Chunks) (bool, error)
	DeleteDocument(ctx context.Context, docID string) (bool, error)
	SearchSimilarDocuments(ctx context.Context, queryEmbed []float32, limit int) ([]models.SimilarDocuments, error)
	UploadChunks(ctx context.Context, chunks []chunks.Chunks) error
	SearchChunkedDocuments(ctx context.Context, embedding []float32, limit int, minScore float32) ([]models.ChunkDetails, error)
	SearchFilteredChunkedDocuments(ctx context.Context, embedding []float32, category string, difficulty string, limit int, minScore float32) ([]models.ChunkDetails, error)
}

func (db *DataBaseStore) Create(ctx context.Context) error {
	if _, err := db.Db.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector"); err != nil {
		return fmt.Errorf("create vector extension: %w", err)
	}

	for tableName, createQuery := range db.Queries.Create {
		_, err := db.Db.Exec(ctx, createQuery)
		if err != nil {
			return err
		} else {
			db.Log.Log.Info().Msgf("Table %s created successfully", tableName)
		}
	}

	return nil
}

func (db *DataBaseStore) SaveDocument(ctx context.Context, document models.Document, embedding []float32) error {
	_, err := db.Db.Exec(ctx, db.Queries.Save.Embedding, document.DocId, document.DocTitle, document.DocContent, document.DocCategory, document.DocSource, document.MetaData, pgvector.NewVector(embedding))
	if err != nil {
		return fmt.Errorf("Error saving embedding data : %w", err)
	}

	return nil
}

func (db *DataBaseStore) UpdateDocument(ctx context.Context, document models.Document, embedding []float32) (bool, error) {
	result, err := db.Db.Exec(ctx, db.Queries.Edit.Document, document.DocTitle, document.DocContent, document.DocCategory, document.DocSource, document.MetaData, pgvector.NewVector(embedding), document.DocId)
	if err != nil {
		return false, fmt.Errorf("error updating document: %w", err)
	}

	return result.RowsAffected() > 0, nil
}

func (db *DataBaseStore) UpdateDocumentWithChunks(ctx context.Context, document models.Document, embedding []float32, documentChunks []chunks.Chunks) (bool, error) {
	dbTxn, err := db.Db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("begin document update transaction: %w", err)
	}
	defer dbTxn.Rollback(ctx)

	result, err := dbTxn.Exec(ctx, db.Queries.Edit.Document, document.DocTitle, document.DocContent, document.DocCategory, document.DocSource, document.MetaData, pgvector.NewVector(embedding), document.DocId)
	if err != nil {
		return false, fmt.Errorf("update document: %w", err)
	}
	if result.RowsAffected() == 0 {
		return false, nil
	}

	if _, err := dbTxn.Exec(ctx, db.Queries.Delete.Chunks, document.DocId); err != nil {
		return false, fmt.Errorf("delete old document chunks: %w", err)
	}

	for _, documentChunk := range documentChunks {
		_, err := dbTxn.Exec(ctx, db.Queries.Save.Chunks, documentChunk.Id, documentChunk.DocId, documentChunk.Content, documentChunk.Index, documentChunk.StartPosition, documentChunk.EndPosition, pgvector.NewVector(documentChunk.ChunkEmbed))
		if err != nil {
			return false, fmt.Errorf("save updated chunk %d: %w", documentChunk.Index, err)
		}
	}

	if _, err := dbTxn.Exec(ctx, db.Queries.Edit.IndexingStatus, document.DocId); err != nil {
		return false, fmt.Errorf("update indexing status: %w", err)
	}
	if err := dbTxn.Commit(ctx); err != nil {
		return false, fmt.Errorf("commit document update transaction: %w", err)
	}

	return true, nil
}

func (db *DataBaseStore) DeleteDocument(ctx context.Context, docID string) (bool, error) {
	dbTxn, err := db.Db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("begin document deletion transaction: %w", err)
	}
	defer dbTxn.Rollback(ctx)

	if _, err := dbTxn.Exec(ctx, db.Queries.Delete.Chunks, docID); err != nil {
		return false, fmt.Errorf("delete document chunks: %w", err)
	}

	result, err := dbTxn.Exec(ctx, db.Queries.Delete.Document, docID)
	if err != nil {
		return false, fmt.Errorf("error deleting document: %w", err)
	}
	if err := dbTxn.Commit(ctx); err != nil {
		return false, fmt.Errorf("commit document deletion transaction: %w", err)
	}

	return result.RowsAffected() > 0, nil
}

func (db *DataBaseStore) SearchSimilarDocuments(ctx context.Context, queryEmbed []float32, limit int) ([]models.SimilarDocuments, error) {
	rows, err := db.Db.Query(ctx, db.Queries.Fetch.SimilarDocuments, pgvector.NewVector(queryEmbed), limit)
	if err != nil {
		return []models.SimilarDocuments{}, err
	}
	defer rows.Close()

	similarDocuments := make([]models.SimilarDocuments, 0)

	for rows.Next() {
		var doc models.SimilarDocuments

		err = rows.Scan(&doc.Id, &doc.DocTitle, &doc.DocContent, &doc.DocCategory, &doc.Similarity)
		if err != nil {
			return []models.SimilarDocuments{}, err
		}

		similarDocuments = append(similarDocuments, doc)
	}
	if err := rows.Err(); err != nil {
		return []models.SimilarDocuments{}, err
	}

	return similarDocuments, nil
}

func (db *DataBaseStore) UploadChunks(ctx context.Context, chunks []chunks.Chunks) error {
	dbTxn, err := db.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin chunk upload transaction: %w", err)
	}
	defer dbTxn.Rollback(ctx)

	docId := ""

	for _, chunk := range chunks {
		// uploading to the database
		_, err := dbTxn.Exec(ctx, db.Queries.Save.Chunks, chunk.Id, chunk.DocId, chunk.Content, chunk.Index, chunk.StartPosition, chunk.EndPosition, pgvector.NewVector(chunk.ChunkEmbed))
		if err != nil {
			return fmt.Errorf("save chunk %d: %w", chunk.Index, err)
		}
		docId = chunk.DocId
	}

	// updating the indexing_status in documents table for the given document
	_, err = dbTxn.Exec(ctx, db.Queries.Edit.IndexingStatus, docId)
	if err != nil {
		return fmt.Errorf("update indexing status: %w", err)
	}
	if err := dbTxn.Commit(ctx); err != nil {
		return fmt.Errorf("commit chunk upload transaction: %w", err)
	}

	return nil
}

func (db *DataBaseStore) SearchChunkedDocuments(ctx context.Context, embedding []float32, limit int, minScore float32) ([]models.ChunkDetails, error) {
	rows, err := db.Db.Query(ctx, db.Queries.Fetch.ChunkedDocuments, pgvector.NewVector(embedding), limit, minScore)
	if err != nil {
		return []models.ChunkDetails{}, err
	}
	defer rows.Close()

	chunkDetails := make([]models.ChunkDetails, 0)

	for rows.Next() {
		var chunk models.ChunkDetails

		err := rows.Scan(&chunk.Id, &chunk.DocumentId, &chunk.Content, &chunk.ChunkIndex, &chunk.Similarity, &chunk.DocumentTitle, &chunk.Category, &chunk.Metadata)
		if err != nil {
			return []models.ChunkDetails{}, err
		}

		chunkDetails = append(chunkDetails, chunk)
	}

	return chunkDetails, nil
}

func (db *DataBaseStore) SearchFilteredChunkedDocuments(ctx context.Context, embedding []float32, category string, difficulty string, limit int, minScore float32) ([]models.ChunkDetails, error) {
	rows, err := db.Db.Query(ctx, db.Queries.Fetch.FilterChunkedDocuments, pgvector.NewVector(embedding), category, difficulty, limit, minScore)
	if err != nil {
		return []models.ChunkDetails{}, err
	}
	defer rows.Close()

	chunkDetails := make([]models.ChunkDetails, 0)

	for rows.Next() {
		var chunk models.ChunkDetails

		err := rows.Scan(&chunk.Id, &chunk.DocumentId, &chunk.Content, &chunk.ChunkIndex, &chunk.Similarity, &chunk.DocumentTitle, &chunk.Category, &chunk.Metadata)
		if err != nil {
			return []models.ChunkDetails{}, err
		}

		chunkDetails = append(chunkDetails, chunk)
	}

	return chunkDetails, nil
}

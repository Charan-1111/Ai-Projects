package database

import (
	"context"
	"fmt"
	"semantic-search/internal/models"

	"github.com/pgvector/pgvector-go"
)

type Repository interface {
	Create(ctx context.Context) error
	SaveDocument(ctx context.Context, document models.Document, embedding []float32) error
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

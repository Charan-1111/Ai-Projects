package server

import (
	"context"
	"fmt"

	"semantic-search/internal/chunks"
	"semantic-search/internal/config"
	"semantic-search/internal/logging"
	"semantic-search/internal/providers"
	"semantic-search/internal/services"
	"semantic-search/internal/store/database"
)

type Application struct {
	log         *logging.Log
	config      *config.Configuration
	llmProvider providers.LLMProvider
	documents   *[]services.Document
	dbStore     database.Repository
	wordChunker *chunks.WordChunker
}

func NewApplication() (*Application, error) {
	log := &logging.Log{}
	log.Initialize()

	applicationConfig := &config.Configuration{}
	if err := applicationConfig.LoadConfig(); err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	databaseStore := &database.DataBaseStore{}
	if err := databaseStore.InitializeDatabaseStore(context.Background(), applicationConfig.Queries, log); err != nil {
		return nil, fmt.Errorf("initialize database store: %w", err)
	}

	geminiClient, err := providers.NewGeminiProvider(context.Background())
	if err != nil {
		return nil, fmt.Errorf("create gemini provider: %w", err)
	}

	wordChunker, err := chunks.NewWordChunker(applicationConfig.ChunkDetails.ChunkSize, applicationConfig.ChunkDetails.ChunkOverlap)
	if err != nil {
		return nil, fmt.Errorf("create word chunker: %w", err)
	}

	return &Application{
		log:         log,
		config:      applicationConfig,
		llmProvider: geminiClient,
		documents:   &[]services.Document{},
		dbStore:     databaseStore,
		wordChunker: wordChunker,
	}, nil
}

func (app *Application) StartServer() error {
	if err := app.dbStore.Create(context.Background()); err != nil {
		return fmt.Errorf("create database tables: %w", err)
	}

	return app.SetupRoutes().Listen(app.config.Address())
}

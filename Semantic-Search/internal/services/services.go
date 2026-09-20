package services

import (
	"semantic-search/internal/chunks"
	"semantic-search/internal/config"
	"semantic-search/internal/logging"
	"semantic-search/internal/providers"
	"semantic-search/internal/store/database"
)

type Service struct {
	config      *config.Configuration
	log         *logging.Log
	llmProvider providers.LLMProvider
	documents   *[]Document
	dbStore     database.Repository
	wordChunker *chunks.WordChunker
}

func NewService(config *config.Configuration, log *logging.Log, llmProvider providers.LLMProvider, documents *[]Document, dbStore database.Repository, wordChunker *chunks.WordChunker) *Service {
	return &Service{
		config:      config,
		log:         log,
		llmProvider: llmProvider,
		documents:   documents,
		dbStore:     dbStore,
		wordChunker: wordChunker,
	}
}

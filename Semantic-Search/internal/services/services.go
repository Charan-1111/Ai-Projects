package services

import (
	"semantic-search/internal/config"
	"semantic-search/internal/logging"
	"semantic-search/internal/providers"
)

type Service struct {
	config      *config.Configuration
	log         *logging.Log
	llmProvider providers.LLMProvider
	documents   *[]Document
}

func NewService(config *config.Configuration, log *logging.Log, llmProvider providers.LLMProvider, documents *[]Document) *Service {
	return &Service{
		config:      config,
		log:         log,
		llmProvider: llmProvider,
		documents:   documents,
	}
}

package services

import (
	"llm-playground/internal/chat"
	"llm-playground/internal/config"
	"llm-playground/internal/provider"
)

type Services struct {
	config              *config.Configuration
	provider            provider.LLMProvider
	inMemoryChatService *chat.InMemoryChatService
}

func NewService(config *config.Configuration, provider provider.LLMProvider, inMemoryChatService *chat.InMemoryChatService) *Services {
	return &Services{
		config:              config,
		provider:            provider,
		inMemoryChatService: inMemoryChatService,
	}
}

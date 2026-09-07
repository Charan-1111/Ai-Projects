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
	persistentChatService *chat.PersistentChatService
}

func NewService(config *config.Configuration, provider provider.LLMProvider, inMemoryChatService *chat.InMemoryChatService, persistentChatService *chat.PersistentChatService) *Services {
	return &Services{
		config:              config,
		provider:            provider,
		inMemoryChatService: inMemoryChatService,
		persistentChatService: persistentChatService,
	}
}

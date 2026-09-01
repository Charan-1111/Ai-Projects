package services

import (
	"llm-playground/internal/config"
	"llm-playground/internal/provider"
)

type Services struct {
	config   *config.Configuration
	provider provider.LLMProvider
}

func NewService(config *config.Configuration, provider provider.LLMProvider) *Services {
	return &Services{
		config:   config,
		provider: provider,
	}
}

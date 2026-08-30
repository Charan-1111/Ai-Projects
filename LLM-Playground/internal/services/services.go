package services

import (
	"llm-playground/internal/config"

	"google.golang.org/genai"
)

type Services struct {
	config *config.Configuration
	client *genai.Client
}

func NewService(config *config.Configuration, client *genai.Client) *Services {
	return &Services{
		config: config,
		client: client,
	}
}

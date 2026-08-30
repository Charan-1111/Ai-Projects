package services

import "llm-playground/internal/config"

type Services struct {
	config *config.Configuration
}

func NewService(config *config.Configuration) *Services {
	return &Services{
		config: config,
	}
}

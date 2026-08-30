package handlers

import (
	"llm-playground/internal/config"
	"llm-playground/internal/services"
)

type Handlers struct {
	config *config.Configuration
	Services *services.Services
}

func NewHandler(config *config.Configuration) *Handlers  {
	return &Handlers{
		config: config,
	}
}

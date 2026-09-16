package handlers

import (
	"semantic-search/internal/config"
	"semantic-search/internal/logging"
	"semantic-search/internal/services"
)

type Handlers struct {
	config  *config.Configuration
	log     *logging.Log
	service *services.Service
}

func NewHandlers(config *config.Configuration, log *logging.Log, service *services.Service) *Handlers {
	return &Handlers{
		config:  config,
		log:     log,
		service: service,
	}
}

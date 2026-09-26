package handlers

import (
	"ask-my-documents/internal/config"
	"ask-my-documents/internal/log"
	"ask-my-documents/internal/services"
)

type Handlers struct {
	config  *config.Configuration
	log     *log.Log
	service *services.Service
}

func NewHandelrs(config *config.Configuration, log *log.Log, service *services.Service) *Handlers {
	return &Handlers{
		config: config,
		log:    log,
		service: service,
	}
}

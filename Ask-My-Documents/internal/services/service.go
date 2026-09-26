package services

import (
	"ask-my-documents/internal/clients"
	"ask-my-documents/internal/config"
	"ask-my-documents/internal/log"
)

type Service struct {
	config     *config.Configuration
	log        *log.Log
	clients    *clients.Clients
	apiFactory clients.ApiFactory
}

func NewService(config *config.Configuration, log *log.Log, clients *clients.Clients, apiFactory clients.ApiFactory) *Service {
	return &Service{
		config:     config,
		log:        log,
		clients:    clients,
		apiFactory: apiFactory,
	}
}

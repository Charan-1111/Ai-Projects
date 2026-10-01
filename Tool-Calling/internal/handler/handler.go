package handler

import (
	"tool-calling/internal/config"
	"tool-calling/internal/services"
)

type Handler struct {
	config   *config.Configuration
	services *services.Service
}

func NewHandler(config *config.Configuration, s *services.Service) *Handler {
	return &Handler{
		config:   config,
		services: s,
	}
}

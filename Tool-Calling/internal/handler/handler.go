package handler

import "tool-calling/internal/config"

type Handler struct {
	config *config.Configuration
}

func NewHandler(config *config.Configuration) *Handler {
	return &Handler{
		config: config,
	}
}

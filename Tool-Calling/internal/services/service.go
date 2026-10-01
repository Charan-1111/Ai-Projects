package services

import "tool-calling/internal/config"

type Service struct {
	config *config.Configuration
}

func NewService(config *config.Configuration) *Service {
	return &Service{
		config: config,
	}
}
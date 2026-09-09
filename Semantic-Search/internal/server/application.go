package server

import (
	"context"
	"fmt"

	"semantic-search/internal/config"
	"semantic-search/internal/logging"
	"semantic-search/internal/providers"
)

type Application struct {
	log        *logging.Log
	config     *config.Configuration
	llmProvider providers.LLMProvider
}

func NewApplication() (*Application, error) {
	log := &logging.Log{}
	log.Initialize()

	applicationConfig := &config.Configuration{}
	if err := applicationConfig.LoadConfig(); err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	geminiClient, err := providers.NewGeminiProvider(context.Background())
	if err != nil {
		return nil, fmt.Errorf("create gemini provider: %w", err)
	}

	return &Application{
		log:        log,
		config:     applicationConfig,
		llmProvider: geminiClient,
	}, nil
}

func (app *Application) StartServer() error {
	return app.SetupRoutes().Listen(app.config.Address())
}

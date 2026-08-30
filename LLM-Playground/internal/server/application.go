package server

import (
	"context"
	"fmt"
	"llm-playground/internal/config"
	"os"

	"google.golang.org/genai"
)

type Application struct {
	config *config.Configuration
	client *genai.Client
}

func NewApplication() (*Application, error) {
	config := &config.Configuration{}

	err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	apiKey := os.Getenv("LLM_PROVIDER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("LLM Provider API key not found")
	}

	client, err := genai.NewClient(
		context.Background(),
		&genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Error initializing genai client : %w", err)
	}

	return &Application{
		config: config,
		client: client,
	}, nil
}

func (app *Application) StartServer() error {
	appServer := app.SetupRoutes()

	err := appServer.Listen(":8000")

	return err
}

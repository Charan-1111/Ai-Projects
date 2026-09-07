package server

import (
	"context"
	"fmt"
	"llm-playground/internal/chat"
	"llm-playground/internal/config"
	"llm-playground/internal/logging"
	"llm-playground/internal/provider"
	"llm-playground/internal/store/database"
	"os"

	"google.golang.org/genai"
)

type Application struct {
	log                 *logging.Log
	config              *config.Configuration
	client              *genai.Client
	provider            provider.LLMProvider
	inMemoryChatService *chat.InMemoryChatService
	dbStore             database.Repository
}

func NewApplication() (*Application, error) {
	log := &logging.Log{}
	log.Initialize()

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

	// llmProvider := provider.GeminiProvider{Client: client}
	llmProvider := provider.NewGeminiProvider(client)

	// creating the inmeory chat se4rvice
	chatService := chat.NewInMemoryChatService(llmProvider)

	databaseStore := &database.DataBaseStore{}
	err = databaseStore.InitializeDatabaseStore(context.Background(), config.Queries, log)
	if err != nil {
		return nil, fmt.Errorf("Error initializing database store : %w", err)
	}

	return &Application{
		log:                 log,
		config:              config,
		client:              client,
		provider:            llmProvider,
		inMemoryChatService: chatService,
		dbStore:             databaseStore,
	}, nil
}

func (app *Application) StartServer() error {
	// creating the tables
	err := app.dbStore.Create(context.Background())
	if err != nil {
		return fmt.Errorf("Error creating tables : %w", err)
	}

	appServer := app.SetupRoutes()

	err = appServer.Listen(":8000")

	return err
}

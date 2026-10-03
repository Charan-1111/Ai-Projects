package server

import (
	"context"
	"fmt"
	"llm-playground/internal/chat"
	"llm-playground/internal/clients"
	"llm-playground/internal/config"
	"llm-playground/internal/logging"
	"llm-playground/internal/provider"
	"llm-playground/internal/store/database"
	"llm-playground/internal/tools"
	"os"

	"google.golang.org/genai"
)

type Application struct {
	log                   *logging.Log
	config                *config.Configuration
	client                *genai.Client
	provider              provider.LLMProvider
	inMemoryChatService   *chat.InMemoryChatService
	persistentChatService *chat.PersistentChatService
	dbStore               database.Repository
	clients               *clients.Clients
	toolRegistry          *tools.ToolRegistry
	apiFactory            clients.ApiFactory
}

func NewApplication() (*Application, error) {
	log := &logging.Log{}
	log.Initialize()

	config := &config.Configuration{}

	err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	databaseStore := &database.DataBaseStore{}
	err = databaseStore.InitializeDatabaseStore(context.Background(), config.Queries, log)
	if err != nil {
		return nil, fmt.Errorf("Error initializing database store : %w", err)
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

	// creating the persistent chat service
	persistentChatService := chat.NewPersistentChatService(llmProvider, databaseStore)

	httpClients := clients.NewClient()
	defaultApiFactory := &clients.DefaultApiFactory{}

	toolRegistry := tools.NewToolRegistry(httpClients, defaultApiFactory, config.Tools)

	return &Application{
		log:                   log,
		config:                config,
		client:                client,
		provider:              llmProvider,
		inMemoryChatService:   chatService,
		persistentChatService: persistentChatService,
		dbStore:               databaseStore,
		clients:               httpClients,
		toolRegistry:          toolRegistry,
		apiFactory:            defaultApiFactory,
	}, nil
}

func (app *Application) StartServer() error {
	// creating the tables
	err := app.dbStore.Create(context.Background())
	if err != nil {
		return fmt.Errorf("Error creating tables : %w", err)
	}

	// fetch the details of the tools...
	app.toolRegistry.RegisterTools()
	registeredToolNames := make(map[string]struct{})
	declarations := app.toolRegistry.GeminiDeclarations()
	for _, declaration := range declarations {
		registeredToolNames[declaration.Name] = struct{}{}
	}
	for modelID, modelConfig := range app.config.AvailableModels {
		for _, requiredTool := range modelConfig.RequiredTools {
			if _, exists := registeredToolNames[requiredTool]; !exists {
				return fmt.Errorf("model %q requires tool %q, but it was not registered", modelID, requiredTool)
			}
		}
	}
	geminiProvider, ok := app.provider.(*provider.GeminiProvider)
	if !ok {
		return fmt.Errorf("configured provider does not support Gemini tool declarations")
	}
	geminiProvider.ConfigureTools(declarations, app.toolRegistry)

	appServer := app.SetupRoutes()

	err = appServer.Listen(":8000")

	return err
}

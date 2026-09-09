package server

import (
	"fmt"

	"semantic-search/internal/config"
)

type Application struct {
	config *config.Configuration
}

func NewApplication() (*Application, error) {
	applicationConfig := &config.Configuration{}
	if err := applicationConfig.LoadConfig(); err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	return &Application{config: applicationConfig}, nil
}

func (app *Application) StartServer() error {
	return app.SetupRoutes().Listen(app.config.Address())
}

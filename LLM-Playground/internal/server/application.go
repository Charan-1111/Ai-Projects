package server

import "llm-playground/internal/config"

type Application struct {
	config *config.Configuration
}

func NewApplication() (*Application, error) {
	config := &config.Configuration{}

	err := config.LoadConfig()
	if err != nil {
		return nil, err
	}


	return &Application{
		config: config,
	}, nil
}

func (app *Application) StartServer() error {
	appServer := app.SetupRoutes()


	err := appServer.Listen(":8000")

	return err
}
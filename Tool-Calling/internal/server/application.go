package server

import "tool-calling/internal/config"

type Application struct {
	config *config.Configuration
}

func NewApplication() (*Application, error) {
	cfg := &config.Configuration{}
	if err := cfg.LoadConfig(); err != nil {
		return nil, err
	}

	return &Application{config: cfg}, nil
}

func (app *Application) StartServer() {
	appServer := app.Router()
	port := app.config.Server.Port
	if port == "" {
		port = ":8003"
	}
	if err := appServer.Listen(port); err != nil {
		panic(err)
	}
}

func (app *Application) StartFiberServer() {
	app.StartServer()
} 
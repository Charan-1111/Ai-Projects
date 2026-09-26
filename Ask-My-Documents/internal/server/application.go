package server

import (
	"ask-my-documents/internal/config"
	"fmt"
)

type Application struct {
	config *config.Configuration
}

func NewApplication() (*Application, error) {
	config := &config.Configuration{}
	err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Loading server configuration : %w", err)
	}

	return &Application{
		config: config,
	}, nil
}

func (app *Application) StartServer() {
	app.StartFiberServer()

}

func (app *Application) StartFiberServer() {
	appServer := app.Router()

	listenErr := make(chan error, 1)

	go func() {
		listenErr <- appServer.Listen(app.config.Server.Port)
	}()



	select {
	case <-listenErr:
		fmt.Println("Server is unable to listen on the port : ", app.config.Server.Port)
	// case 
	// TODO : Need to have the graceful shutdown code ready
	}
}

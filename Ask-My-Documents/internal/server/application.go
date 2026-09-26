package server

import (
	"ask-my-documents/internal/clients"
	"ask-my-documents/internal/config"
	"ask-my-documents/internal/log"
	"fmt"
)

type Application struct {
	config     *config.Configuration
	log        *log.Log
	clients    *clients.Clients
	apiFactory clients.ApiFactory
}

func NewApplication() (*Application, error) {
	log := &log.Log{}
	log.Initialize()

	config := &config.Configuration{}
	err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Loading server configuration : %w", err)
	}

	defaultFactory := &clients.DefaultApiFactory{}

	httpClients := clients.NewClient()

	

	return &Application{
		config:  config,
		log:     log,
		clients: httpClients,
		apiFactory: defaultFactory,
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
		app.log.Log.Error().Err(<-listenErr).Msg("Server is unable to listen on the port : " + app.config.Server.Port)
		// case
		// TODO : Need to have the graceful shutdown code ready
	}
}

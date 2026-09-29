package server

import "tool-calling/internal/config"

type Application struct {
	config *config.Configuration
}

func NewApplication() (*Application, error) {
	config := &config.Configuration{}
	err := config.LoadConfig()
	if err != nil {

	}

	return &Application{
		config: config,
	}, nil
}

func (app *Application) StartServer() {

}

func (app *Application) StartFiberServer() {
	// appServer :=- 
} 
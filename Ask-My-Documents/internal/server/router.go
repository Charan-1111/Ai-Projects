package server

import (
	"ask-my-documents/internal/handlers"
	"ask-my-documents/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func (app *Application) Router() *fiber.App {
	appServer := fiber.New()

	// need to have the middlewares
	appServer.Use(cors.New())
	// appServer.Use(recover.New())
	// need to have a request id middleware as well
	
	service := services.NewService(app.config, app.log, app.clients, app.apiFactory)
	handler := handlers.NewHandelrs(app.config, app.log, service)

	apiGroup := appServer.Group("/rag")
	apiGroup.Post("/ask/documents", handler.AskDocuments)
	return appServer
}

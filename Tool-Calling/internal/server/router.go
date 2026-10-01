package server

import (
	"tool-calling/internal/handler"
	"tool-calling/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func (app *Application) Router() *fiber.App {
	appServer := fiber.New()

	// add middlewares
	appServer.Use(cors.New())
	// appServer.Use(recover.New())

	apiGroup := appServer.Group("/assistant")

	services := services.NewService(app.config)
	handler := handler.NewHandler(app.config, services)

	apiGroup.Post("/chat", handler.GetTime)

	toolsGroup := appServer.Group("/tools")
	toolsGroup.Get("/", handler.ListTools)
	toolsGroup.Post("/execute", handler.ExecuteTool)
	return appServer
}

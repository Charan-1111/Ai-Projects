package server

import (
	"tool-calling/internal/handler"
	"tool-calling/internal/mcpServer"
	"tool-calling/internal/services"

	"github.com/gofiber/fiber/v3/middleware/adaptor"
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

	appServer.Post("/get-time", handler.GetCurrentTime)

	// MCP
	appServer.All("/mcp", adaptor.HTTPHandler(mcpServer.NewHandler(services)))

	toolsGroup := appServer.Group("/tools")
	toolsGroup.Get("/", handler.ListTools)
	toolsGroup.Post("/execute", handler.ExecuteTool)
	return appServer
}

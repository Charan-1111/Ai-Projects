package server

import (
	"llm-playground/internal/handlers"
	"llm-playground/internal/middleware"
	"llm-playground/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/sse"
)

func (app *Application) SetupRoutes() *fiber.App {
	appServer := fiber.New()

	appServer.Use(recover.New())
	//adding request id middleware
	appServer.Use(middleware.RequestId)

	handler := handlers.NewHandler(app.config)
	handler.Services = services.NewService(app.config, app.provider)

	appServer.Get("/health", handler.HealthCheck)

	apiGroup := appServer.Group("/v1/llm")
	apiGroup.Get("/models/available", handler.AvailableModels)
	apiGroup.Post("/generate", handler.GenerateResponse)
	apiGroup.Post("/generate/stream", sse.New(sse.Config{Handler: handler.GenerateStreamResponse}))
	return appServer
}

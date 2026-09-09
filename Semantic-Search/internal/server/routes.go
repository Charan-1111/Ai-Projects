package server

import (
	"semantic-search/internal/handlers"
	"semantic-search/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func (app *Application) SetupRoutes() *fiber.App {
	appServer := fiber.New()
	appServer.Use(recover.New())
	appServer.Use(cors.New())
	appServer.Use(middleware.RequestID)

	appHandlers := handlers.NewHandlers()
	appServer.Get("/health", appHandlers.HealthCheck)

	return appServer
}

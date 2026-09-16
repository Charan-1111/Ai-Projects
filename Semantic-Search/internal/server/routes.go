package server

import (
	"semantic-search/internal/handlers"
	"semantic-search/internal/middleware"
	"semantic-search/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func (app *Application) SetupRoutes() *fiber.App {
	appServer := fiber.New(fiber.Config{})

	appServer.Use(recover.New())
	appServer.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
	}))
	appServer.Use(middleware.RequestID)

	service := services.NewService(app.config, app.log, app.llmProvider, app.documents)
	appHandlers := handlers.NewHandlers(app.config, app.log, service)

	appServer.Get("/health", appHandlers.HealthCheck)

	apiGroup := appServer.Group("/semantic")
	apiGroup.Post("/embed", appHandlers.EmbedText)
	apiGroup.Post("/compare", appHandlers.CompareSimilarity)

	docGroup := apiGroup.Group("/document")
	docGroup.Post("/inject", appHandlers.InjectDocument)

	return appServer
}

package server

import (
	"tool-calling/internal/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func (app *Application) Router() *fiber.App {
	appServer := fiber.New()

	// add middlewares
	appServer.Use(cors.New())
	// appServer.Use(recover.New())

	apiGroup := appServer.Group("/assistant")

	handler := handler.NewHandler(app.config)

	apiGroup.Post("/chat", handler.GetTime)
	return appServer
}

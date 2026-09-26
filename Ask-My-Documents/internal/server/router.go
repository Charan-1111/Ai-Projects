package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func (ap *Application) Router() *fiber.App {
	appServer := fiber.New()

	// need to have the middlewares
	appServer.Use(cors.New())
	// appServer.Use(recover.New())
	// need to have a request id middleware as well

	return appServer
}

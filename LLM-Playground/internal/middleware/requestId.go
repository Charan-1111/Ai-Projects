package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func RequestId(c fiber.Ctx) error {
	requestId := c.Get("X-Request-Id")
	if requestId == "" {
		requestId = uuid.NewString()
	}

	c.Set("X-Request-Id", requestId)
	c.Locals("requestId", requestId)

	return c.Next()
}

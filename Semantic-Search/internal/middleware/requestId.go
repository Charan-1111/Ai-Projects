package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func RequestID(c fiber.Ctx) error {
	requestID := c.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.NewString()
	}

	c.Set("X-Request-ID", requestID)
	c.Locals("requestId", requestID)

	return c.Next()
}

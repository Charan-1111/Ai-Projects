package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func SessionId(c fiber.Ctx) error {
	sessionId := c.Get("X-Session-Id")
	if sessionId == "" {
		sessionId = uuid.NewString()
	}

	c.Set("X-Session-Id", sessionId)
	c.Locals("sessionId", sessionId)

	return c.Next()
}

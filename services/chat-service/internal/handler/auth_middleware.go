package handler

import (
	"unbound-v2/services/chat-service/internal/grpcclient"

	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(client *grpcclient.AuthClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := c.Get("Authorization")
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}

		userID, err := client.ValidateToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}

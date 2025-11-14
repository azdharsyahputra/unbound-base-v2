package handler

import (
	"strings"
	"unbound-v2/services/chat-service/internal/grpcclient"

	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(client *grpcclient.AuthClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := ""

		// 1. Coba ambil dari header
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				token = parts[1]
			}
		}

		// 2. Fallback: ambil dari query param `token`
		if token == "" {
			token = c.Query("token")
		}

		// 3. Kalau masih kosong → unauthorized
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}

		// 4. Validate via Auth-Service gRPC
		userID, err := client.ValidateToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}

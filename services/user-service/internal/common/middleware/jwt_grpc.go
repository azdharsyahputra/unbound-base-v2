package middleware

import (
	"strings"

	"unbound-v2/services/user-service/internal/grpcclient"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected(authClient *grpcclient.AuthClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing Authorization header"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		res, err := authClient.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token validation failed"})
		}
		if !res.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("userID", res.UserId)
		return c.Next()
	}
}

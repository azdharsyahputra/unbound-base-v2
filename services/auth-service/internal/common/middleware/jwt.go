package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"unbound-v2/services/auth-service/internal/auth"
)

func JWTProtected(authSvc *auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		uid, err := authSvc.ParseToken(tokenStr)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		c.Locals("userID", uid)
		return c.Next()
	}
}

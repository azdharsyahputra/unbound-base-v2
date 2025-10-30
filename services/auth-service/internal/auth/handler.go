package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, svc *AuthService) {
	r := app.Group("/auth")

	r.Post("/register", func(c *fiber.Ctx) error {
		var req RegisterReq
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		u, err := svc.Register(req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
		})
	})

	r.Post("/login", func(c *fiber.Ctx) error {
		var req LoginReq
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}

		tok, err := svc.Login(req)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return c.JSON(fiber.Map{
			"success":        true,
			"access_token":   tok.AccessToken,
			"refresh_token":  tok.RefreshToken,
			"token_type":     "Bearer",
			"expires_in_sec": 86400,
		})
	})

	r.Post("/refresh", func(c *fiber.Ctx) error {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}

		tok, err := svc.RefreshAccess(body.RefreshToken)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return c.JSON(fiber.Map{
			"success":        true,
			"access_token":   tok.AccessToken,
			"refresh_token":  tok.RefreshToken,
			"token_type":     "Bearer",
			"expires_in_sec": 86400,
		})
	})

	r.Post("/logout", func(c *fiber.Ctx) error {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&body); err != nil || body.RefreshToken == "" {
			return fiber.NewError(fiber.StatusBadRequest, "refresh_token required")
		}

		if err := svc.DB.Where("token = ?", body.RefreshToken).Delete(&RefreshToken{}).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to logout")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Logged out successfully",
		})
	})
}

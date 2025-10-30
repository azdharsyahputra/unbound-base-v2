package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	r := app.Group("/users")
	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Users endpoint (use /users/:username)"})
	})
}

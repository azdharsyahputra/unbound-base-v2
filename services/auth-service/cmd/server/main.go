package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

    "unbound-v2/services/auth-service/internal/auth"
    "unbound-v2/services/auth-service/internal/common/db"
    "unbound-v2/services/auth-service/internal/common/middleware"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	app := fiber.New()
	app.Use(middleware.JSONResponseMiddleware)

	// Database
	database := db.Connect()

	// Auth service & route
	authSvc := auth.NewAuthService(database)
	auth.RegisterRoutes(app, database, authSvc)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "auth-service",
			"status":  "running ✅",
		})
	})

	log.Fatal(app.Listen(":8081")) // port khusus auth-service
}

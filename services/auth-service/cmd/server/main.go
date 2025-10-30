package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env kalau ada
	_ = godotenv.Load()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "auth-service",
			"status":  "running ✅",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 auth-service running on port %s", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

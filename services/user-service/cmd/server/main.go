package main

import (
	"fmt"
	"log"
	"os"

	"unbound-v2/services/user-service/internal/common/db"
	"unbound-v2/services/user-service/internal/grpcclient"
	"unbound-v2/services/user-service/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()
	database := db.Connect()

	authAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authAddr == "" {
		authAddr = "localhost:50051"
	}
	authClient := grpcclient.NewAuthClient(authAddr)

	user.RegisterRoutes(app, database)
	user.RegisterProfileRoutes(app, database)
	user.RegisterFollowRoutes(app, database, authClient)

	// ✅ Tambahkan /check di sini
	app.Get("/check", func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(400).JSON(fiber.Map{"error": "missing Authorization header"})
		}

		res, err := authClient.ValidateToken(token)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if !res.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		return c.JSON(fiber.Map{
			"message": "token valid ✅",
			"user_id": res.UserId,
		})
	})

	// ✅ Healthcheck root
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "user-service",
			"status":  "running ✅",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("🚀 user-service running on port %s", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

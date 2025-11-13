package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"unbound-v2/notification-service/internal/db"
	"unbound-v2/notification-service/internal/events"
	"unbound-v2/notification-service/internal/notification"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()

	// DB connect
	database := db.Connect()
	database.AutoMigrate(&notification.Notification{})
	log.Println("Notification DB connected & migrated")

	// gRPC Auth client
	authAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authAddr == "" {
		authAddr = "auth-service:50051"
	}

	// Kafka consumer
	go events.StartKafkaConsumer(database)

	// Routes
	notification.RegisterNotificationRoutes(app, database)

	// Root healthcheck
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "notification-service",
			"status":  "running",
		})
	})

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("🚀 notification-service running on %s", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

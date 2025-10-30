package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"unbound-v2/services/auth-service/internal/auth"
	"unbound-v2/services/auth-service/internal/common/db"
	"unbound-v2/services/auth-service/internal/common/middleware"
	"unbound-v2/services/auth-service/internal/grpcserver"
	authpb "unbound-v2/shared/proto/auth"
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

	// ✅ Tambahkan route root untuk healthcheck
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "auth-service",
			"status":  "running ✅",
		})
	})

	// ✅ Jalankan gRPC server di goroutine terpisah
	go func() {
		grpcPort := os.Getenv("AUTH_GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "50051"
		}

		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
		if err != nil {
			log.Fatalf("❌ failed to listen: %v", err)
		}

		s := grpc.NewServer()
		authpb.RegisterAuthServiceServer(s, grpcserver.NewAuthGRPCServer(authSvc))
		log.Printf("🛰 gRPC auth-service running on port %s", grpcPort)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("❌ failed to serve gRPC: %v", err)
		}
	}()

	// 🚀 Jalankan HTTP server (Fiber)
	log.Println("🚀 HTTP auth-service running on port 8081")
	log.Fatal(app.Listen(":8081"))
}

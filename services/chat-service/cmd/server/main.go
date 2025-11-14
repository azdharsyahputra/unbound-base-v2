package main

import (
	"log"

	"unbound-v2/services/chat-service/internal/config"
	"unbound-v2/services/chat-service/internal/grpcclient"
	"unbound-v2/services/chat-service/internal/handler"
	"unbound-v2/services/chat-service/internal/migration"
	"unbound-v2/services/chat-service/internal/repository"
	"unbound-v2/services/chat-service/internal/route"
	"unbound-v2/services/chat-service/internal/service"
	ws "unbound-v2/services/chat-service/internal/websocket"

	authpb "unbound-v2/shared/proto/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ======================================================
// DATABASE INIT
// ======================================================

func initDB(cfg *config.Config) *gorm.DB {
	dsn := cfg.DBDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect DB: %v", err)
	}
	log.Println("📦 Connected to Postgres")
	return db
}

// ======================================================
// KAFKA WRITER INIT
// ======================================================

func initKafka(cfg *config.Config) *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{cfg.Kafka.Broker},
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	})
	log.Println("📡 Kafka writer ready:", cfg.Kafka.Broker)
	return writer
}

// ======================================================
// GRPC CONNECTION INIT (AUTH ONLY)
// ======================================================

func initAuthConn(cfg *config.Config) *grpc.ClientConn {
	conn, err := grpc.Dial(cfg.AuthServiceURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed connect Auth gRPC: %v", err)
	}
	log.Println("🔐 Connected to Auth-Service gRPC:", cfg.AuthServiceURL)
	return conn
}

// ======================================================
// MAIN FUNCTION
// ======================================================

func main() {
	// LOAD CONFIG
	cfg := config.Load()
	log.Println("🚀 Starting Chat-Service on port:", cfg.App.Port)

	// INIT DB
	db := initDB(cfg)

	// MIGRATION
	migration.Migrate(db)

	// INIT KAFKA
	kafkaWriter := initKafka(cfg)

	// INIT AUTH gRPC CLIENT
	authConn := initAuthConn(cfg)
	authGrpc := authpb.NewAuthServiceClient(authConn)
	authClient := grpcclient.NewAuthClient(authGrpc) // wrapper

	// INIT REPOSITORIES
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// INIT EVENT SERVICE
	eventSvc := service.NewEventService(kafkaWriter)

	// INIT CORE SERVICES
	chatSvc := service.NewChatService(chatRepo, authClient)
	messageSvc := service.NewMessageService(messageRepo, chatRepo, eventSvc)

	// INIT WEBSOCKET HUB & SERVICE
	hub := ws.NewHub()
	wsSvc := service.NewWebSocketService(hub, messageSvc, chatSvc)

	// INIT HANDLERS
	chatHandler := handler.NewChatHandler(chatSvc, messageSvc)
	wsHandler := handler.NewWebSocketHandler(hub, wsSvc, authClient)

	// AUTH MIDDLEWARE (VALIDATE JWT via AUTH-SERVICE)
	authMiddleware := handler.NewAuthMiddleware(authClient)

	// FIBER APP
	app := fiber.New()

	// ROUTES
	route.RegisterRoutes(app, chatHandler, wsHandler, authMiddleware, hub)

	// START SERVER
	log.Println("✅ Chat-Service running on port:", cfg.App.Port)

	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatalf("❌ Failed starting Chat-Service: %v", err)
	}
}

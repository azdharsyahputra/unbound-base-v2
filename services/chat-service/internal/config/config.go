package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Port string
	}

	DB struct {
		Host string
		User string
		Pass string
		Name string
		Port string
	}

	Kafka struct {
		Broker string
		Topic  string
	}

	AuthServiceURL string
	UserServiceURL string
}

func Load() *Config {
	// Load .env automatically
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found, using system env")
	}

	cfg := &Config{}

	cfg.App.Port = os.Getenv("PORT")

	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Pass = os.Getenv("DB_PASS")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.Port = os.Getenv("DB_PORT")

	cfg.Kafka.Broker = os.Getenv("KAFKA_BROKER")
	cfg.Kafka.Topic = os.Getenv("CHAT_TOPIC")

	cfg.AuthServiceURL = os.Getenv("AUTH_GRPC_ADDR")
	cfg.UserServiceURL = os.Getenv("USER_GRPC_ADDR")

	return cfg
}

func (c *Config) DBDSN() string {
	return "host=" + c.DB.Host +
		" user=" + c.DB.User +
		" password=" + c.DB.Pass +
		" dbname=" + c.DB.Name +
		" port=" + c.DB.Port +
		" sslmode=disable"
}

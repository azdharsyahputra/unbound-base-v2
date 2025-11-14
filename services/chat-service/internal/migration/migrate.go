package migration

import (
	"log"

	"unbound-v2/services/chat-service/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.Chat{},
		&model.Message{},
	)

	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("database migrated successfully")
}

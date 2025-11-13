package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"

	"unbound-v2/notification-service/internal/notification"
)

type FollowEvent struct {
	Type        string `json:"type"`
	FollowerID  uint   `json:"follower_id"`
	FollowingID uint   `json:"following_id"`
}

func StartKafkaConsumer(db *gorm.DB) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"redpanda:9092"},
		Topic:   "notifications",
		GroupID: "notification-service-group",
	})

	log.Println("📥 Notification service listening Kafka topic: notifications")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("❌ Kafka read error:", err)
			continue
		}

		var event FollowEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Println("❌ Failed to decode event:", err)
			continue
		}

		// hanya FOLLOW event
		if event.Type != "follow" {
			continue
		}

		message := "Someone followed you"

		notif := notification.Notification{
			UserID:   event.FollowingID,
			SenderID: event.FollowerID,
			Type:     "follow",
			Message:  message,
		}

		if err := db.Create(&notif).Error; err != nil {
			log.Println("❌ Failed to store notification:", err)
		} else {
			log.Printf("🔔 Follow notification stored (to user=%d from=%d)\n",
				event.FollowingID, event.FollowerID)
		}
	}

}

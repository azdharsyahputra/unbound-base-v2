package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer() *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("redpanda:9092"),
		Topic:    "notifications",
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) PublishFollowEvent(eventType string, followerID, followingID uint64) error {
	event := map[string]any{
		"type":         eventType,
		"follower_id":  followerID,
		"following_id": followingID,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.writer.WriteMessages(context.Background(), kafka.Message{Value: data})
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("📤 Published %s event: follower=%d → following=%d", eventType, followerID, followingID)
	return nil
}

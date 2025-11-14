package service

import (
	"context"
	"encoding/json"
	"log"
	"unbound-v2/services/chat-service/internal/model"

	"github.com/segmentio/kafka-go"
)

type EventService struct {
	Writer *kafka.Writer
}

func NewEventService(writer *kafka.Writer) *EventService {
	return &EventService{Writer: writer}
}

func (e *EventService) publish(eventType string, payload interface{}) {
	data, _ := json.Marshal(map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	})

	err := e.Writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: data,
		},
	)

	if err != nil {
		log.Println("❌ Failed to publish event:", err)
	}
}

func (e *EventService) MessageCreated(msg model.Message) {
	e.publish("message_created", msg)
}

func (e *EventService) MessageDelivered(chatID uint, userID uint) {
	e.publish("message_delivered", map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	})
}

func (e *EventService) MessageRead(chatID uint, userID uint) {
	e.publish("message_read", map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	})
}

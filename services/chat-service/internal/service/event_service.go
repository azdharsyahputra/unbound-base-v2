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

func (e *EventService) publish(eventKey string, payload interface{}) {
	data, _ := json.Marshal(payload)

	err := e.Writer.WriteMessages(
		context.Background(), // wajib!
		kafka.Message{
			Key:   []byte(eventKey),
			Value: data,
		},
	)

	if err != nil {
		log.Printf("Failed to publish event %s: %v", eventKey, err)
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

package service

import (
	"time"
	"unbound-v2/services/chat-service/internal/model"
	"unbound-v2/services/chat-service/internal/repository"
)

type MessageService struct {
	MessageRepo *repository.MessageRepository
	ChatRepo    *repository.ChatRepository
	EventSvc    EventPublisher // kirim Kafka event
}

type EventPublisher interface {
	MessageCreated(msg model.Message)
	MessageDelivered(chatID uint, userID uint)
	MessageRead(chatID uint, userID uint)
}

func NewMessageService(msgRepo *repository.MessageRepository, chatRepo *repository.ChatRepository, eventSvc EventPublisher) *MessageService {
	return &MessageService{
		MessageRepo: msgRepo,
		ChatRepo:    chatRepo,
		EventSvc:    eventSvc,
	}
}

// Retrieve messages
func (s *MessageService) GetMessages(chatID uint) ([]model.Message, error) {
	return s.MessageRepo.ListByChatID(chatID)
}

// Send new message
func (s *MessageService) SendMessage(chatID, senderID uint, content string) (*model.Message, error) {

	msg := &model.Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
		Status:   "sent",
		IsRead:   false,
	}

	// save to database
	if err := s.MessageRepo.Create(msg); err != nil {
		return nil, err
	}

	// publish event for notification-service
	s.EventSvc.MessageCreated(*msg)

	return msg, nil
}

// Mark messages as delivered
func (s *MessageService) MarkDelivered(chatID, userID uint) error {
	err := s.MessageRepo.MarkDelivered(chatID, userID)
	if err == nil {
		s.EventSvc.MessageDelivered(chatID, userID)
	}
	return err
}

// Mark messages as read
func (s *MessageService) MarkRead(chatID, userID uint) error {
	now := time.Now()

	err := s.MessageRepo.MarkRead(chatID, userID, now)
	if err == nil {
		s.EventSvc.MessageRead(chatID, userID)
	}
	return err
}

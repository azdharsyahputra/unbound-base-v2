package service

import (
	"unbound-v2/services/chat-service/internal/model"
	"unbound-v2/services/chat-service/internal/repository"
)

type MessageService struct {
	MessageRepo *repository.MessageRepository
	ChatRepo    *repository.ChatRepository
	EventSvc    EventPublisher
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

// SEND message
func (s *MessageService) SendMessage(chatID, senderID uint, content string) (*model.Message, error) {

	msg := &model.Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
		Status:   "sent",
		IsRead:   false,
	}

	if err := s.MessageRepo.Create(msg); err != nil {
		return nil, err
	}

	s.EventSvc.MessageCreated(*msg)

	return msg, nil
}

// DELIVERED
func (s *MessageService) MarkDelivered(chatID, userID uint) error {

	err := s.MessageRepo.MarkDelivered(chatID, userID)
	if err == nil {
		s.EventSvc.MessageDelivered(chatID, userID)
	}

	return err
}

// READ RECEIPT
func (s *MessageService) MarkAsRead(msgID, userID uint) error {

	err := s.MessageRepo.MarkAsRead(msgID, userID)
	if err == nil {
		s.EventSvc.MessageRead(msgID, userID)
	}

	return err
}

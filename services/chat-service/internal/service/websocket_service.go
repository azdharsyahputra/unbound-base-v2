package service

import (
	"time"
	"unbound-v2/services/chat-service/internal/model"
	"unbound-v2/services/chat-service/internal/websocket"
)

type WebSocketService struct {
	Hub        *websocket.Hub
	MessageSvc *MessageService
	ChatSvc    *ChatService
}

func NewWebSocketService(
	hub *websocket.Hub,
	messageSvc *MessageService,
	chatSvc *ChatService,
) *WebSocketService {
	return &WebSocketService{
		Hub:        hub,
		MessageSvc: messageSvc,
		ChatSvc:    chatSvc,
	}
}

// Handle message send from websocket incoming event
func (s *WebSocketService) HandleIncomingMessage(chatID, senderID uint, content string) (*model.Message, error) {
	// process message with MessageService
	msg, err := s.MessageSvc.SendMessage(chatID, senderID, content)
	if err != nil {
		return nil, err
	}

	// broadcast through websocket hub
	s.Hub.Broadcast(websocket.Payload{
		Type:      "message",
		ChatID:    chatID,
		MessageID: msg.ID,
		SenderID:  senderID,
		Content:   msg.Content,
		Status:    msg.Status,
		Timestamp: msg.CreatedAt,
	})

	return msg, nil
}

// When user opens websocket connection, mark messages as delivered
func (s *WebSocketService) HandleDelivered(chatID, userID uint) error {
	err := s.MessageSvc.MarkDelivered(chatID, userID)
	if err == nil {
		s.Hub.Broadcast(websocket.Payload{
			Type:      "status_update",
			ChatID:    chatID,
			SenderID:  userID,
			Status:    "delivered",
			Timestamp: time.Now(),
		})
	}
	return err
}

package service

import (
	"time"
	"unbound-v2/services/chat-service/internal/model"
	ws "unbound-v2/services/chat-service/internal/websocket"
)

type WebSocketService struct {
	Hub        *ws.Hub
	MessageSvc *MessageService
	ChatSvc    *ChatService
}

func NewWebSocketService(
	hub *ws.Hub,
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

	msg, err := s.MessageSvc.SendMessage(chatID, senderID, content)
	if err != nil {
		return nil, err
	}

	s.Hub.Broadcast(ws.Payload{
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
		s.Hub.Broadcast(ws.Payload{
			Type:      "status_update",
			ChatID:    chatID,
			SenderID:  userID,
			Status:    "delivered",
			Timestamp: time.Now(),
		})
	}

	return err
}

// READ RECEIPT
func (s *WebSocketService) HandleRead(chatID, userID, msgID uint) error {

	err := s.MessageSvc.MarkAsRead(msgID, userID)
	if err != nil {
		return err
	}

	s.Hub.BroadcastToRoom(chatID, map[string]interface{}{
		"type":       "status_update",
		"chat_id":    chatID,
		"message_id": msgID,
		"sender_id":  userID,
		"status":     "read",
		"timestamp":  time.Now().UTC(),
	})

	return nil
}

// TYPING indicator
func (s *WebSocketService) HandleTyping(chatID, userID uint) {
	s.Hub.Broadcast(ws.Payload{
		Type:      "typing",
		ChatID:    chatID,
		SenderID:  userID,
		Timestamp: time.Now(),
	})
}

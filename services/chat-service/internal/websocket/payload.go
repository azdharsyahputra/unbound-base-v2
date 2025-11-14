package websocket

import "time"

type Payload struct {
	Type      string    `json:"type"`
	ChatID    uint      `json:"chat_id"`
	MessageID uint      `json:"message_id,omitempty"`
	SenderID  uint      `json:"sender_id"`
	Content   string    `json:"content,omitempty"`
	Status    string    `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

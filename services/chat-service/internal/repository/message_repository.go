package repository

import (
	"time"
	"unbound-v2/services/chat-service/internal/model"

	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

// Get all messages in a chat
func (r *MessageRepository) ListByChatID(chatID uint) ([]model.Message, error) {
	var messages []model.Message
	err := r.DB.
		Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Find(&messages).Error

	return messages, err
}

// Create new message
func (r *MessageRepository) Create(msg *model.Message) error {
	return r.DB.Create(msg).Error
}

// Mark messages as delivered
func (r *MessageRepository) MarkDelivered(chatID, userID uint) error {
	return r.DB.
		Model(&model.Message{}).
		Where("chat_id = ? AND sender_id != ? AND status = ?", chatID, userID, "sent").
		Update("status", "delivered").Error
}

// Mark messages as read
func (r *MessageRepository) MarkRead(chatID, userID uint, t time.Time) error {
	return r.DB.
		Model(&model.Message{}).
		Where("chat_id = ? AND sender_id != ? AND status != ?", chatID, userID, "read").
		Updates(map[string]interface{}{
			"status":  "read",
			"is_read": true,
			"read_at": t,
		}).Error
}

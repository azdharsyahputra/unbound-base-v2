package repository

import (
	"unbound-v2/services/chat-service/internal/model"

	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

// LIST all messages in a chat (urut naik by created_at)
func (r *MessageRepository) ListByChatID(chatID uint) ([]model.Message, error) {
	var msgs []model.Message

	err := r.DB.
		Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Find(&msgs).Error

	return msgs, err
}

// CREATE message
func (r *MessageRepository) Create(msg *model.Message) error {
	return r.DB.Create(msg).Error
}

// MARK delivered
func (r *MessageRepository) MarkDelivered(chatID, userID uint) error {
	return r.DB.
		Model(&model.Message{}).
		Where("chat_id = ? AND sender_id != ? AND status = ?", chatID, userID, "sent").
		Update("status", "delivered").
		Error
}

// MARK read
func (r *MessageRepository) MarkAsRead(msgID, userID uint) error {
	return r.DB.
		Model(&model.Message{}).
		Where("id = ? AND sender_id != ?", msgID, userID).
		Updates(map[string]interface{}{
			"status":  "read",
			"is_read": true,
		}).
		Error
}

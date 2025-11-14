package repository

import (
	"unbound-v2/services/chat-service/internal/model"

	"gorm.io/gorm"
)

type ChatRepository struct {
	DB *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{DB: db}
}

// Cari chat yang sudah ada antara 2 user
func (r *ChatRepository) FindBetweenUsers(user1, user2 uint) (*model.Chat, error) {
	var chat model.Chat
	err := r.DB.
		Where(`(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)`,
			user1, user2, user2, user1).
		First(&chat).Error

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// Create chat baru
func (r *ChatRepository) Create(chat *model.Chat) error {
	return r.DB.Create(chat).Error
}

// Get chat dengan preload messages (optional)
func (r *ChatRepository) GetByID(id uint) (*model.Chat, error) {
	var chat model.Chat
	err := r.DB.
		Where("id = ?", id).
		First(&chat).Error

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// List chat untuk 1 user
func (r *ChatRepository) ListByUser(userID uint) ([]model.Chat, error) {
	var chats []model.Chat
	err := r.DB.
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Find(&chats).Error

	return chats, err
}

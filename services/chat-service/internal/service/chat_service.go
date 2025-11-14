package service

import (
	"errors"

	"unbound-v2/services/chat-service/internal/model"
	"unbound-v2/services/chat-service/internal/repository"
)

var ErrUserNotFound = errors.New("target user not found")

type ChatService struct {
	ChatRepo   *repository.ChatRepository
	UserClient UserClient // gRPC ke user-service
}

type UserClient interface {
	VerifyUserExists(userID uint) (bool, error)
}

func NewChatService(chatRepo *repository.ChatRepository, userClient UserClient) *ChatService {
	return &ChatService{
		ChatRepo:   chatRepo,
		UserClient: userClient,
	}
}

// Get existing chat or create new one
func (s *ChatService) GetOrCreateChat(user1ID, user2ID uint) (*model.Chat, error) {

	// cek apakah target user exist di user-service
	exist, err := s.UserClient.VerifyUserExists(user2ID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrUserNotFound
	}

	// cek chat yang sudah ada
	chat, err := s.ChatRepo.FindBetweenUsers(user1ID, user2ID)
	if err == nil {
		return chat, nil
	}

	// kalau tidak ada -> create baru
	newChat := &model.Chat{
		User1ID: user1ID,
		User2ID: user2ID,
	}

	if err := s.ChatRepo.Create(newChat); err != nil {
		return nil, err
	}

	return newChat, nil
}

func (s *ChatService) GetByID(chatID uint) (*model.Chat, error) {
	return s.ChatRepo.GetByID(chatID)
}

func (s *ChatService) ListByUser(userID uint) ([]model.Chat, error) {
	return s.ChatRepo.ListByUser(userID)
}

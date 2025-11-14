package service

import (
	"errors"
	"fmt"

	"unbound-v2/services/chat-service/internal/model"
	"unbound-v2/services/chat-service/internal/repository"
)

var (
	ErrUserNotFound = errors.New("target user not found")
	ErrSelfChat     = errors.New("cannot create chat with yourself")
)

type ChatService struct {
	ChatRepo       *repository.ChatRepository
	AuthUserClient AuthUserClient // gRPC ke AUTH-SERVICE
}

type AuthUserClient interface {
	VerifyUserExists(userID uint) (bool, error)
}

func NewChatService(chatRepo *repository.ChatRepository, authClient AuthUserClient) *ChatService {
	return &ChatService{
		ChatRepo:       chatRepo,
		AuthUserClient: authClient,
	}
}

func (s *ChatService) GetOrCreateChat(user1ID, user2ID uint) (*model.Chat, error) {

	// =====================
	// 1. Prevent self-chat
	// =====================
	if user1ID == user2ID {
		return nil, ErrSelfChat
	}

	// =====================================
	// 2. Normalize ordering (smaller first)
	// supaya chat 2-5 == chat 5-2
	// =====================================
	u1 := user1ID
	u2 := user2ID

	if u1 > u2 {
		u1, u2 = u2, u1
	}

	// =============================
	// 3. Verify user target exists
	// =============================
	exist, err := s.AuthUserClient.VerifyUserExists(user2ID)
	if err != nil {
		return nil, fmt.Errorf("user verification failed: %w", err)
	}
	if !exist {
		return nil, ErrUserNotFound
	}

	// =============================
	// 4. Find existing chat
	// =============================
	chat, err := s.ChatRepo.FindBetweenUsers(u1, u2)
	if err == nil {
		return chat, nil
	}

	// =============================
	// 5. Create new chat
	// =============================
	newChat := &model.Chat{
		User1ID: u1,
		User2ID: u2,
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

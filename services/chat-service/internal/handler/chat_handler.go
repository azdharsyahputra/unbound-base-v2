package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"unbound-v2/services/chat-service/internal/service"
)

type ChatHandler struct {
	ChatSvc    *service.ChatService
	MessageSvc *service.MessageService
}

func NewChatHandler(chatSvc *service.ChatService, msgSvc *service.MessageService) *ChatHandler {
	return &ChatHandler{
		ChatSvc:    chatSvc,
		MessageSvc: msgSvc,
	}
}

// GET /chats
func (h *ChatHandler) ListChats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	chats, err := h.ChatSvc.ListByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(chats)
}

// POST /chats/:user_id
func (h *ChatHandler) GetOrCreateChat(c *fiber.Ctx) error {
	userID1 := c.Locals("user_id").(uint)

	user2Str := c.Params("user_id")
	userID2, err := strconv.ParseUint(user2Str, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user ID")
	}

	chat, err := h.ChatSvc.GetOrCreateChat(userID1, uint(userID2))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(chat)
}

// GET /chats/:chat_id/messages
func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
	chatIDStr := c.Params("chat_id")
	chatID64, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid chat_id")
	}
	chatID := uint(chatID64)

	// sekarang langsung pakai repository ListByChatID()
	messages, err := h.MessageSvc.MessageRepo.ListByChatID(chatID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(messages)
}

// POST /chats/:chat_id/messages
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	chatIDStr := c.Params("chat_id")
	chatID64, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid chat_id")
	}
	chatID := uint(chatID64)

	var body struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&body); err != nil || body.Content == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid message body")
	}

	msg, err := h.MessageSvc.SendMessage(chatID, userID, body.Content)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(msg)
}

// PUT /chats/:chat_id/read
func (h *ChatHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req struct {
		MessageID uint `json:"message_id"`
	}

	if err := c.BodyParser(&req); err != nil || req.MessageID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	err := h.MessageSvc.MarkAsRead(req.MessageID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "marked as read",
	})
}

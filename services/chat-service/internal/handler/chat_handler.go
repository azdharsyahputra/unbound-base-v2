package handler

import (
	"strconv"
	"unbound-v2/services/chat-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	ChatSvc    *service.ChatService
	MessageSvc *service.MessageService
}

func NewChatHandler(chatSvc *service.ChatService, messageSvc *service.MessageService) *ChatHandler {
	return &ChatHandler{
		ChatSvc:    chatSvc,
		MessageSvc: messageSvc,
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
	userID := c.Locals("user_id").(uint)

	targetStr := c.Params("user_id")
	targetID64, err := strconv.ParseUint(targetStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user_id")
	}
	targetID := uint(targetID64)

	chat, err := h.ChatSvc.GetOrCreateChat(userID, targetID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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

	messages, err := h.MessageSvc.GetMessages(chatID)
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

	var req struct {
		Content string `json:"content"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON")
	}

	msg, err := h.MessageSvc.SendMessage(chatID, userID, req.Content)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(msg)
}

// PUT /chats/:chat_id/read
func (h *ChatHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	chatIDStr := c.Params("chat_id")
	chatID64, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid chat_id")
	}
	chatID := uint(chatID64)

	if err := h.MessageSvc.MarkRead(chatID, userID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"status": "ok"})
}

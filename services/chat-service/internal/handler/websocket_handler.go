package handler

import (
	"encoding/json"
	"strconv"

	"unbound-v2/services/chat-service/internal/grpcclient"
	"unbound-v2/services/chat-service/internal/service"
	ws "unbound-v2/services/chat-service/internal/websocket"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebSocketHandler struct {
	Hub        *ws.Hub
	WSService  *service.WebSocketService
	AuthClient *grpcclient.AuthClient
}

func NewWebSocketHandler(
	hub *ws.Hub,
	wsService *service.WebSocketService,
	authClient *grpcclient.AuthClient,
) *WebSocketHandler {
	return &WebSocketHandler{
		Hub:        hub,
		WSService:  wsService,
		AuthClient: authClient,
	}
}

// /ws/chat/:chat_id?token=xxxx
func (h *WebSocketHandler) Handle(c *fiber.Ctx) error {

	chatIDStr := c.Params("chat_id")
	chatID64, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid chat_id")
	}
	chatID := uint(chatID64)

	token := c.Query("token")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing token")
	}

	// validate token via Auth-Service gRPC
	userID, err := h.AuthClient.ValidateToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	return websocket.New(func(conn *websocket.Conn) {

		// REGISTER CONNECTION
		h.Hub.HandleConnection(
			chatID,
			userID,
			conn,

			// --------------------------
			// ON MESSAGE RECEIVED
			// --------------------------
			func(chatID, userID uint, data []byte) {

				var incoming struct {
					Content string `json:"content"`
					Read    uint   `json:"read"`
					Typing  bool   `json:"typing"`
				}

				if err := json.Unmarshal(data, &incoming); err != nil {
					return
				}

				// TYPING INDICATOR
				if incoming.Typing {
					h.WSService.HandleTyping(chatID, userID)
					return
				}

				// READ RECEIPT
				if incoming.Read != 0 {
					h.WSService.HandleRead(chatID, userID, incoming.Read)
					return
				}

				// SEND MESSAGE
				if incoming.Content != "" {
					h.WSService.HandleIncomingMessage(chatID, userID, incoming.Content)
					return
				}
			},

			// --------------------------
			// ON CONNECT
			// --------------------------
			func() {
				h.WSService.HandleDelivered(chatID, userID)
			},

			// ON DISCONNECT
			nil,
		)

	})(c)
}

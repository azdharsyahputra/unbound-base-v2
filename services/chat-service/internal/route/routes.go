package route

import (
	"unbound-v2/services/chat-service/internal/handler"
	ws "unbound-v2/services/chat-service/internal/websocket"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	chatHandler *handler.ChatHandler,
	wsHandler *handler.WebSocketHandler,
	authMiddleware fiber.Handler,
	hub *ws.Hub,
) {

	// ===== REST ROUTES =====
	r := app.Group("/chats", authMiddleware)

	r.Get("/", chatHandler.ListChats)
	r.Post("/:user_id", chatHandler.GetOrCreateChat)
	r.Get("/:chat_id/messages", chatHandler.GetMessages)
	r.Post("/:chat_id/messages", chatHandler.SendMessage)
	r.Put("/:chat_id/read", chatHandler.MarkAsRead)

	// ===== WEBSOCKET ROUTE =====
	// Step 1: auth + upgrade validation
	app.Get("/ws/chat/:chat_id", authMiddleware, func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return wsHandler.Handle(c) // ← LANGSUNG KESINI
		}
		return fiber.ErrUpgradeRequired
	})
}

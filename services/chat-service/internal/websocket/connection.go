package websocket

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

type ConnectionHandler func(chatID uint, userID uint, msg []byte)

func (h *Hub) HandleConnection(chatID uint, userID uint, conn *websocket.Conn, onMessage ConnectionHandler, onConnect func(), onClose func()) {

	h.JoinRoom(chatID, userID, conn)
	log.Printf("User %d connected to chat %d", userID, chatID)

	if onConnect != nil {
		onConnect()
	}

	defer func() {
		h.LeaveRoom(chatID, conn)
		log.Printf("User %d disconnected from chat %d", userID, chatID)

		if onClose != nil {
			onClose()
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if onMessage != nil {
			onMessage(chatID, userID, data)
		}
	}
}

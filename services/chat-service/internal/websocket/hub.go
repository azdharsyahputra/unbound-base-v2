package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	Rooms map[uint]map[*websocket.Conn]uint
	Mutex sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[uint]map[*websocket.Conn]uint),
	}
}

// Register connection into room
func (h *Hub) JoinRoom(chatID uint, userID uint, conn *websocket.Conn) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	if _, exists := h.Rooms[chatID]; !exists {
		h.Rooms[chatID] = make(map[*websocket.Conn]uint)
	}

	h.Rooms[chatID][conn] = userID
}

// Remove connection
func (h *Hub) LeaveRoom(chatID uint, conn *websocket.Conn) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	delete(h.Rooms[chatID], conn)
	conn.Close()
}

// Broadcast to specific chat room
func (h *Hub) Broadcast(payload Payload) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	conns := h.Rooms[payload.ChatID]

	data, _ := json.Marshal(payload)

	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("WebSocket write error:", err)
			conn.Close()
			delete(conns, conn)
		}
	}
}

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

// Broadcast to specific chat room (Payload based)
func (h *Hub) Broadcast(payload Payload) {
	h.Mutex.Lock()
	conns := h.Rooms[payload.ChatID]
	h.Mutex.Unlock()

	data, _ := json.Marshal(payload)

	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("WebSocket write error:", err)
			conn.Close()

			h.Mutex.Lock()
			delete(conns, conn)
			h.Mutex.Unlock()
		}
	}
}

func (h *Hub) BroadcastToRoom(chatID uint, data interface{}) {
	h.Mutex.Lock()
	conns := h.Rooms[chatID]
	h.Mutex.Unlock()

	if conns == nil {
		return
	}

	bytes, _ := json.Marshal(data)

	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
			log.Println("WebSocket write error:", err)
			conn.Close()

			h.Mutex.Lock()
			delete(conns, conn)
			h.Mutex.Unlock()
		}
	}
}

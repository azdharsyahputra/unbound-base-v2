package websocket

type Message struct {
	UserID  uint
	Payload any
}

type Hub struct {
	clients    map[uint]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client.UserID][client]; ok {
				delete(h.clients[client.UserID], client)
				close(client.send)
			}

		case msg := <-h.broadcast:
			for client := range h.clients[msg.UserID] {
				client.send <- msg.Payload
			}
		}
	}
}

func (h *Hub) Broadcast(msg Message) {
	h.broadcast <- msg
}

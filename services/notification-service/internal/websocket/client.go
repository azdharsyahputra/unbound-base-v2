package websocket

import (
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Hub    *Hub
	send   chan any
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// Client tidak mengirim apa-apa, ini receive-only
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}

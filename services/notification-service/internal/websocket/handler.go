package websocket

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RegisterWebsocketRoute(app *fiber.App, hub *Hub) {
	app.Get("/ws/:userID", websocket.New(func(conn *websocket.Conn) {
		idStr := conn.Params("userID")
		uid, _ := strconv.ParseUint(idStr, 10, 64)

		client := &Client{
			UserID: uint(uid),
			Conn:   conn,
			Hub:    hub,
			send:   make(chan any),
		}

		// Register client
		hub.register <- client

		go client.WritePump()
		client.ReadPump()
	}))
}

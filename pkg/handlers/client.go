package handlers

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	Hub        *Hub
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		connection: conn,
		Hub:        hub,
	}
}

func (c *Client) ReadMessages() {
	for {
		messageType, p, err := c.connection.ReadMessage()
		if err != nil {
			slog.Error("Error reading message", slog.Any("error", err))
			return
		}
		slog.Info("Received message from client",
			slog.Group("message", slog.Int("message_type", messageType), slog.String("message", string(p))),
		)
	}
}

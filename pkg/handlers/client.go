package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/gorilla/websocket"
	"github.com/grez-lucas/go-websocket-drawing/pkg/dto"
)

type ClientList map[*Client]bool

type Client struct {
	connection   *websocket.Conn
	Hub          *Hub
	messageQueue chan ([]byte)
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		connection:   conn,
		Hub:          hub,
		messageQueue: make(chan ([]byte)),
	}
}

func (c *Client) ReadMessages() {
	defer c.Hub.removeClient(c) // If something goes wrong, disconnect the client
	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Error reading message", slog.Any("error", err))
			}
			break
		}
		slog.Info("Received message from client",
			slog.Group("message", slog.Int("message_type", messageType), slog.String("message", string(payload))),
		)

		// TODO: Unmarshal the message and send it to all clients

		var statusMessage dto.DrawDTO
		if err := json.Unmarshal(payload, &statusMessage); err != nil {
			slog.Error("Error unmasrhalling message", slog.String("message", string(payload)), slog.Any("error", err))
			break
		}
		for client := range c.Hub.clients {
			client.messageQueue <- payload
		}
	}
}

func (c *Client) WriteMessages() {
	defer func() {
		c.Hub.removeClient(c)
	}()

	for {
		msg, ok := <-c.messageQueue
		if !ok {
			if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
				slog.Warn("connection closed", slog.Any("error", err))
			}
			return
		}
		if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
			slog.Error("Failed to send message", slog.Any("error", err))
		}

		slog.Info("Message sent")
	}
}

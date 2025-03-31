package ws

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
	MessageQueue chan ([]byte)
	Chatroom     string
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		connection:   conn,
		Hub:          hub,
		MessageQueue: make(chan ([]byte)),
		Chatroom:     "1",
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

		// Route the Message
		var msg dto.Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			slog.Error("failed to unmarshal DTO message", slog.Any("error", err))
			continue
		}
		c.Hub.routeMessage(msg, c)
	}
}

func (c *Client) WriteMessages() {
	defer func() {
		c.Hub.removeClient(c)
	}()

	for {
		msg, ok := <-c.MessageQueue
		if !ok {
			if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
				slog.Warn("connection closed", slog.Any("error", err))
			}
			return
		}

		// Marshal the message, it must follow DrawDTO struct
		if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
			slog.Error("Failed to send message", slog.Any("error", err))
		}

		slog.Info("Message sent")
	}
}

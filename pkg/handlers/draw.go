package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/grez-lucas/go-websocket-drawing/pkg/dto"
	"github.com/grez-lucas/go-websocket-drawing/pkg/ws"
)

type DrawHandler struct{}

func NewDrawHandler() *DrawHandler {
	return &DrawHandler{}
}

func (h *DrawHandler) Handle(message dto.Message, c *ws.Client) error {
	slog.Info("Handling Draw Message")
	var drawDTOMsg dto.DrawDTOMessage
	if err := json.Unmarshal(message.Payload, &drawDTOMsg); err != nil {
		return fmt.Errorf("failed to unmarshal drawDTOMessage: %w", err)
	}

	drawDTOBytes, err := json.Marshal(drawDTOMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal drawDTOMessage: %w", err)
	}
	// Broadcast the message to all clients in the hub in that chatroom
	slog.Debug("Drawing message to all clients in chatroom", slog.String("chatroom", c.Chatroom))
	for client := range c.Hub.Clients {
		if client.Chatroom == c.Chatroom {
			client.MessageQueue <- drawDTOBytes
		}
	}
	return nil
}

package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/grez-lucas/go-websocket-drawing/pkg/dto"
	"github.com/grez-lucas/go-websocket-drawing/pkg/ws"
)

type ChatRoomHandler struct{}

func NewChatRoomHandler() *ChatRoomHandler {
	return &ChatRoomHandler{}
}

func (h *ChatRoomHandler) Handle(message dto.Message, c *ws.Client) error {
	// Marshal payload into wanted format
	var changeRoomMsg dto.ChangeRoomMessage
	if err := json.Unmarshal(message.Payload, &changeRoomMsg); err != nil {
		return fmt.Errorf("failed to unmarshal change room message: %w", err)
	}

	c.Chatroom = changeRoomMsg.Name
	return nil
}

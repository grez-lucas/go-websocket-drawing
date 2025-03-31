package dto

import "encoding/json"

const (
	MessageDraw       = "draw_message"
	MessageChangeRoom = "change_room"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type DrawDTOMessage struct {
	StartX int `json:"startX"`
	StartY int `json:"startY"`
	EndX   int `json:"endX"`
	EndY   int `json:"endY"`
}

type ChangeRoomMessage struct {
	Name string `json:"name"`
}

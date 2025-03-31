package main

import (
	"log/slog"

	"github.com/grez-lucas/go-websocket-drawing/pkg/api"
	"github.com/grez-lucas/go-websocket-drawing/pkg/handlers"
	"github.com/grez-lucas/go-websocket-drawing/pkg/ws"
)

func main() {
	slog.Info("Hello world!")

	chatRoomHandler := handlers.NewChatRoomHandler()
	drawHandler := handlers.NewDrawHandler()

	wsHub := ws.NewHub(chatRoomHandler, drawHandler)

	wsServer := api.NewWSServer(":8000", wsHub)
	wsServer.Init()
}

package main

import (
	"log/slog"

	"github.com/grez-lucas/go-websocket-drawing/pkg/api"
	"github.com/grez-lucas/go-websocket-drawing/pkg/handlers"
)

func main() {
	slog.Info("Hello world!")

	wsHub := handlers.NewHub()

	wsServer := api.NewWSServer(":8000", wsHub)
	wsServer.Init()
}

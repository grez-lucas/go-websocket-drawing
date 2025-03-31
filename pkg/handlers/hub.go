package handlers

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Hub struct {
	clients ClientList
	sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(ClientList),
	}
}

func (h *Hub) Upgrade(w http.ResponseWriter, r *http.Request) {
	slog.Info("New connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Error upgrading to WS connection", slog.Any("error", err))
	}

	client := NewClient(conn, h)
	h.addClient(client)

	go client.ReadMessages()
}

func (h *Hub) addClient(c *Client) {
	h.Lock()
	defer h.Unlock()

	if isConnected, ok := h.clients[c]; !ok {
		h.clients[c] = true
		if !isConnected {
			slog.Info("New client added to client list", slog.Any("client", c.connection.LocalAddr()))
			return
		}
		slog.Warn("Client is already registered to WS Hub", slog.Any("client", c.connection.LocalAddr()))
		return
	}
	slog.Warn("Registered client connection added to hub list", slog.Any("client", c.connection.LocalAddr()))
}

func (h *Hub) removeClient(c *Client) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.clients[c]; !ok {
		slog.Error("Could not remove client from hub list, unregistered", slog.Any("client", c))
		return
	}

	c.connection.Close()
	delete(h.clients, c)
}

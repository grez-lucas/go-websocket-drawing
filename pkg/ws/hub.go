package ws

import (
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/grez-lucas/go-websocket-drawing/pkg/dto"
)

var ErrMessageNotSupported = errors.New("this message type is not supported")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Hub struct {
	Clients ClientList
	sync.RWMutex
	handlers map[string]IMessageHandler
}

func NewHub(changeRoomHandler, drawHandler IMessageHandler) *Hub {
	handlers := map[string]IMessageHandler{
		dto.MessageChangeRoom: changeRoomHandler,
		dto.MessageDraw:       drawHandler,
	}
	return &Hub{
		Clients:  make(ClientList),
		handlers: handlers,
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
	go client.WriteMessages()
}

func (h *Hub) routeMessage(msg dto.Message, client *Client) error {
	handler, ok := h.handlers[msg.Type]
	if !ok {
		return ErrMessageNotSupported
	}

	if err := handler.Handle(msg, client); err != nil {
		return err
	}

	return nil
}

func (h *Hub) addClient(c *Client) {
	h.Lock()
	defer h.Unlock()

	if isConnected, ok := h.Clients[c]; !ok {
		h.Clients[c] = true
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

	if _, ok := h.Clients[c]; !ok {
		slog.Error("Could not remove client from hub list, unregistered", slog.Any("client", c))
		return
	}

	c.connection.Close()
	delete(h.Clients, c)
}

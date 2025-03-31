package ws

import (
	"net/http"

	"github.com/grez-lucas/go-websocket-drawing/pkg/dto"
)

type IWsUpgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request)
}

type IMessageHandler interface {
	Handle(dto.Message, *Client) error
}

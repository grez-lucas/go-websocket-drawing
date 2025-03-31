package api

import (
	"log"
	"net/http"

	"github.com/grez-lucas/go-websocket-drawing/pkg/ws"
)

type WSServer struct {
	listenAddr string
	wsUprader  ws.IWsUpgrader
}

func NewWSServer(listenAddr string, upgrader ws.IWsUpgrader) *WSServer {
	return &WSServer{
		listenAddr: listenAddr,
		wsUprader:  upgrader,
	}
}

func (wss *WSServer) Init() {
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", wss.wsUprader.Upgrade)
	log.Fatal(http.ListenAndServe(wss.listenAddr, nil))
	// router := http.NewServeMux()
}

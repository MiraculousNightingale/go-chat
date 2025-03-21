package http

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type option func(*httpServer)

func defaultHttpServer() *httpServer {
	return &httpServer{
		host:       "0.0.0.0",
		port:       "8080",
		wsUpgrader: websocket.Upgrader{},
		sessions:   make(map[string]*session),
		log:        slog.Default(),
	}
}

// Default host is 0.0.0.0
func WithHost(host string) option {
	return func(hs *httpServer) {
		hs.host = host
	}
}

// Default port is 8080
func WithPort(port string) option {
	return func(hs *httpServer) {
		hs.port = port
	}
}

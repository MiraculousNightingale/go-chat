package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

type httpServer struct {
	host       string
	port       string
	wsUpgrader websocket.Upgrader
	sessions   map[string]*session
	log        *slog.Logger
}

func NewServer(options ...option) *httpServer {
	server := defaultHttpServer()

	for _, opt := range options {
		opt(server)
	}

	return server
}

func (srv *httpServer) getAddress() string {
	return fmt.Sprintf("%s:%s", srv.host, srv.port)
}

func (srv *httpServer) DebugPrint() {
	l := slog.Default()

	l.Info("Running HTTP server with the following configuration...")
	l.Info("Host: " + srv.host)
	l.Info("Port: " + srv.port)
}

// Runs the server and blocks current goroutine
func (srv *httpServer) Start() {
	http.HandleFunc("/ws", srv.funcWS)

	addr := srv.getAddress()

	srv.log.Info("Listening and serving", "address", addr)
	http.ListenAndServe(
		addr,
		nil,
	)
}

func (srv *httpServer) Stop() {
	for id, session := range srv.sessions {
		session.Stop()
		delete(srv.sessions, id)
	}

	srv.log.Info("Stopped")
}

func (srv *httpServer) funcWS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("HTTP method must be 'GET'"))

		srv.log.Debug("Received non-GET method in request header", "method", r.Method)

		return
	}

	conn, err := srv.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		srv.log.Error("Failed to upgrade", "error", err)

		return
	}

	session, id := NewSession(srv, conn)
	srv.sessions[id] = session

	session.Start()
}

func (srv *httpServer) deleteSession(id string) {
	delete(srv.sessions, id)
}

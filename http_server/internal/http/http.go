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

func New(options ...option) *httpServer {
	server := defaultHttpServer()

	for _, opt := range options {
		opt(server)
	}

	return server
}

func (srv *httpServer) DebugPrint() {
	fmt.Println("HTTP server debug print:")
	fmt.Printf("host: %s\n", srv.host)
	fmt.Printf("port: %s\n", srv.port)
}

// Runs the server and blocks current goroutine
func (srv *httpServer) Start() {
	http.HandleFunc("/ws", srv.funcWS)

	http.ListenAndServe(
		fmt.Sprintf("%s:%s", srv.host, srv.port),
		nil,
	)
}

func (srv *httpServer) Stop() {
	for id, session := range srv.sessions {
		session.Stop()
		delete(srv.sessions, id)
	}
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

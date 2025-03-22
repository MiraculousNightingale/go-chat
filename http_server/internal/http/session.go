package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

type session struct {
	srv  *httpServer
	id   string
	conn *websocket.Conn
	log  *slog.Logger
}

func NewSession(srv *httpServer, conn *websocket.Conn) (*session, string) {
	id := xid.New().String()
	log := slog.Default().With("component", "session", "id", id)

	return &session{
		srv:  srv,
		id:   id,
		conn: conn,
		log:  log,
	}, id
}

func (s *session) Start() {
	wm, err := WrapMessage(WelcomeMessage{
		Message: "Welcome to go-chat!",
		Date:    time.Now().Format(time.DateTime),
		Version: "0.0.1",
	})
	if err != nil {
		s.log.Error("Failed to wrap welcome message", "error", err)
		s.Stop()

		return
	}

	if err := s.conn.WriteJSON(wm); err != nil {
		s.log.Error("Failed to write welcome message", "error", err)
		s.Stop()

		return
	}

	go func() {
		defer s.Stop()
		defer s.log.Debug("Exited handling loop")

		for {
			_, r, err := s.conn.NextReader()
			if err != nil {
				s.log.Error("Failed to get next reader", "error", err)

				return
			}

			var envelope MessageEnvelope
			if err := json.NewDecoder(r).Decode(&envelope); err != nil {
				s.log.Error("Failed to decode envelope", "error", err)

				continue
			}

			if err = s.conn.WriteMessage(websocket.TextMessage, fmt.Appendf(nil, "received %s message", envelope.Type)); err != nil {
				s.log.Error("Failed to write response", "error", err)

				continue
			}
		}
	}()

	s.log.Debug("Started")
}

func (s *session) Stop() {

	s.conn.Close()
	s.srv.deleteSession(s.id)

	s.log.Debug("Stopped")
}

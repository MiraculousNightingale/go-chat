package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func startServer() (*httpServer, string) {
	srv := NewServer()

	go func() {
		srv.Start()
	}()

	url := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/ws",
	}

	return srv, url.String()
}

func connect(t *testing.T, url string) (*websocket.Conn, *http.Response) {
	conn, res, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		t.Fatal(err)
	}

	return conn, res
}

func TestConnect(t *testing.T) {
	srv, url := startServer()
	defer srv.Stop()

	conn, res := connect(t, url)
	defer conn.Close()

	assert.Equal(t, http.StatusSwitchingProtocols, res.StatusCode, "Connection must upgrade to websocket")

	var env MessageEnvelope
	if err := conn.ReadJSON(&env); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, env.ID, "ID must be 0")
	assert.Equal(t, "welcome", env.Type, "Message type must be 'welcome'")

	var msg WelcomeMessage
	if err := json.Unmarshal(env.Payload, &msg); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Welcome to go-chat!", msg.Message)
	assert.Equal(t, "0.0.1", msg.Version)
}

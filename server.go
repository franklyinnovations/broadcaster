package broadcaster

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// A Server is the main class of this package, pass it to http.Handle on a
// chosen path to start a broadcast server.
type Server struct {
	// Invoked upon initial connection, can be used to enforce access control.
	CanConnect func(data map[string]string) bool

	// Invoked upon channel subscription, can be used to enforce access control
	// for channels.
	CanSubscribe func(data map[string]string, channel string) bool

	// Can be used to configure buffer sizes etc.
	// See http://godoc.org/github.com/gorilla/websocket#Upgrader
	Upgrader websocket.Upgrader

	hub      hub
	prepared bool
}

type clientMessage struct {
	Id      string            `json:"id,omitempty"`
	Type    string            `json:"type"`
	Channel string            `json:"channel,omitempty"`
	Data    map[string]string `json:"data,omitempty"`
}

func (s *Server) Prepare() error {
	s.hub.Server = s
	go s.hub.Run()
	s.prepared = true
	return nil
}

// Main HTTP server.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.prepared {
		http.Error(w, "Prepare() not called on broadcaster.Server", 500)
		return
	}

	if r.Method == "GET" {
		s.handleWebsocket(w, r)
	} else if r.Method == "POST" {
		s.handleLongPoll(w, r)
	}
}

func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	// Always a new client, easy!
	newWebsocketClient(w, r, s)
}

func (s *Server) handleLongPoll(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) Stats() (Stats, error) {
	return s.hub.Stats()
}
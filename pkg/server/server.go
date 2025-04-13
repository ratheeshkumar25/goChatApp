package server

import (
	"github.com/gorilla/mux"
	"github.com/ratheeshkumar25/chatApp/pkg/api"
)

// Server represents the HTTP server
type Server struct {
	Router  *mux.Router
	Handler *api.Handler
}

// NewServer creates a new server with routes configured
func NewServer(handler *api.Handler) *Server {
	server := &Server{
		Router:  mux.NewRouter(),
		Handler: handler,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	s.Router.HandleFunc("/join", s.Handler.JoinHandler).Methods("GET")
	s.Router.HandleFunc("/leave", s.Handler.LeaveHandler).Methods("GET")
	s.Router.HandleFunc("/send", s.Handler.SendHandler).Methods("GET")
	s.Router.HandleFunc("/messages", s.Handler.MessagesHandler).Methods("GET")
}

package api

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	handlers *Handlers
	router   *mux.Router
}

func NewServer(handlers *Handlers) *Server {
	router := mux.NewRouter()

	router.Use(RequestLogger)

	return &Server{
		handlers: handlers,
		router:   router,
	}
}

func (s *Server) RegisterHandlers() {
	s.router.Path("/health").Methods("GET").HandlerFunc(s.handlers.HandleHealthCheck)

	s.router.Path("/tasks").Methods("POST").HandlerFunc(s.handlers.HandleCreateTask)
	s.router.Path("/tasks/{id}").Methods("GET").HandlerFunc(s.handlers.HandleGetTask)
	s.router.Path("/tasks").Methods("GET").HandlerFunc(s.handlers.HandleListTasks)
	s.router.Path("/tasks/{id}").Methods("PUT").HandlerFunc(s.handlers.HandleUpdateTask)
	s.router.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(s.handlers.HandleDeleteTask)
}

func (s *Server) Start() error {
	if err := http.ListenAndServe(":9091", s.router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

package api

import (
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	handlers *Handlers
	router   *mux.Router
}

func NewServer(handlers *Handlers) *Server {
	router := mux.NewRouter()

	router.Use(RequestLogger)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteJSONError(w, errors.New("route not found"), http.StatusNotFound)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if methods := allowedMethods(router, r); len(methods) > 0 {
			w.Header().Set("Allow", strings.Join(methods, ", "))
		}
		WriteJSONError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
	})

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

func allowedMethods(router *mux.Router, r *http.Request) []string {
	var methods []string
	seen := make(map[string]struct{})
	_ = router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		routeMethods, err := route.GetMethods()
		if err != nil {
			return nil
		}

		for _, method := range routeMethods {
			request := r.Clone(r.Context())
			request.Method = method
			var match mux.RouteMatch
			if route.Match(request, &match) {
				if _, exists := seen[method]; !exists {
					methods = append(methods, method)
					seen[method] = struct{}{}
				}
			}
		}
		return nil
	})

	sort.Strings(methods)
	return methods
}

func (s *Server) Start() error {
	if err := http.ListenAndServe(":9091", s.router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

package http

import (
	"errors"
	"net/http"

	"github.com/dapr/go-sdk/service/common"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// NewService creates new Service.
func NewService(address string) common.Service {
	return newServer(address, nil)
}

// NewServiceWithMux creates new Service with existing http mux.
func NewServiceWithMux(address string, mux *mux.Router) common.Service {
	return newServer(address, mux)
}

func newServer(address string, router *mux.Router) *Server {
	if router == nil {
		router = mux.NewRouter()
	}
	return &Server{
		address:            address,
		mux:                router,
		topicSubscriptions: make([]*common.Subscription, 0),
	}
}

// Server is the HTTP server wrapping mux many Dapr helpers.
type Server struct {
	address            string
	mux                *mux.Router
	topicSubscriptions []*common.Subscription
}

// Start starts the HTTP handler. Blocks while serving.
func (s *Server) Start() error {
	s.registerSubscribeHandler()
	c := negroni.New()
	c.UseHandler(s.mux)
	return errors.Unwrap(http.ListenAndServe(s.address, c))
}

// Stop stops previously started HTTP service.
func (s *Server) Stop() error {
	// TODO: implement service stop
	return nil
}

func setOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	w.Header().Set("Allow", "POST,OPTIONS")
}

func optionsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			setOptions(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

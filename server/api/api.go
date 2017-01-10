package api

import (
	"net/http"

	"github.com/ian-kent/gofigure"
	"github.com/ian-kent/service.go/log"
	"github.com/mailhog/mh2/server"
	"github.com/mailhog/mh2/server/api/backend"
	"github.com/mailhog/mh2/server/api/handlers"

	mh2http "github.com/mailhog/mh2/server/http"
)

type apiServer struct {
	config     Config
	httpServer mh2http.Server
	backend    backend.API
	exit       bool
}

// NewServer returns a new server
func NewServer() (server.Server, error) {
	var apiConfig = Config{
		BindAddr: "0.0.0.0:8025",
		Backend:  "mongodb",
	}

	if err := gofigure.Gofigure(&apiConfig); err != nil {
		return nil, err
	}

	be, err := backend.New(apiConfig.Backend)
	if err != nil {
		return nil, err
	}

	api := &apiServer{
		config:     apiConfig,
		httpServer: mh2http.Get(apiConfig.BindAddr),
		backend:    be,
	}

	api.httpServer.Router().Get("/api/healthcheck", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})

	api.httpServer.Router().Get("/messages", handlers.List{API: api.backend}.ServeHTTP)

	return api, nil
}

// Start starts the server
func (s *apiServer) Start() error {
	log.Debug("api: starting server", log.Data{
		"bind_addr": s.config.BindAddr,
	})

	return s.httpServer.Start()
}

// Stop stops the server
func (s *apiServer) Stop() error {
	log.Debug("api: stopping server", nil)
	s.exit = true
	s.backend.Close()
	return s.httpServer.Stop()
}

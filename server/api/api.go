package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ian-kent/gofigure"
	"github.com/ian-kent/service.go/log"
	"github.com/mailhog/mh2/backend/inmemory"
	"github.com/mailhog/mh2/server"
	"github.com/mailhog/mh2/server/api/backend"
	"github.com/mailhog/mh2/server/api/handlers"

	mh2http "github.com/mailhog/mh2/server/http"
)

type apiServer struct {
	config     Config
	httpServer mh2http.Server
	apiBackend backend.API
	exit       bool
}

// NewServer returns a new server
func NewServer() (server.Server, error) {
	var apiConfig = Config{
		BindAddr: "0.0.0.0:8025",
		Backend:  "inmemory",
	}

	if err := gofigure.Gofigure(&apiConfig); err != nil {
		return nil, err
	}

	api := &apiServer{
		config:     apiConfig,
		httpServer: mh2http.Get(apiConfig.BindAddr),
	}

	api.httpServer.Router().Get("/api/healthcheck", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})

	api.httpServer.Router().Get("/messages", handlers.List{API: api.apiBackend}.ServeHTTP)

	return api, nil
}

// Start starts the server
func (s *apiServer) Start() error {
	log.Debug("api: starting server", log.Data{
		"bind_addr": s.config.BindAddr,
	})

	// TODO: refactor this so it's a registration/lookup not a switch statement

	switch strings.ToLower(s.config.Backend) {
	case "inmemory":
		inmem, err := inmemory.New()
		if err != nil {
			return err
		}
		s.apiBackend = inmem
	default:
		return errors.New("api: unrecognised message receiver type")
	}

	return s.httpServer.Start()
}

// Stop stops the server
func (s *apiServer) Stop() error {
	log.Debug("api: stopping server", nil)
	s.exit = true
	s.apiBackend.Close()
	return s.httpServer.Stop()
}

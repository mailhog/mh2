package ui

import (
	"net/http"

	"github.com/ian-kent/service.go/log"

	"github.com/ian-kent/gofigure"
	"github.com/mailhog/mh2/server"

	mh2http "github.com/mailhog/mh2/server/http"
)

// Server is an SMTP server
type uiServer struct {
	config     Config
	httpServer mh2http.Server
}

// NewServer returns a new server
func NewServer() (server.Server, error) {
	var uiConfig = Config{
		BindAddr: "0.0.0.0:8025",
	}

	if err := gofigure.Gofigure(&uiConfig); err != nil {
		return nil, err
	}

	ui := &uiServer{
		config:     uiConfig,
		httpServer: mh2http.Get(uiConfig.BindAddr),
	}

	ui.httpServer.Router().Get("/healthcheck", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})

	return ui, nil
}

// Start starts the server
func (s *uiServer) Start() error {
	log.Debug("ui: starting server", log.Data{
		"bind_addr": s.config.BindAddr,
	})
	return s.httpServer.Start()
}

func (s *uiServer) Stop() error {
	log.Debug("ui: stopping server", nil)
	return s.httpServer.Stop()
}

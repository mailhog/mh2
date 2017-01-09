package http

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/pat"
)

var servers = make(map[string]Server)

var mutex sync.Mutex

// Server represents a HTTP server
type Server interface {
	Router() *pat.Router
	Start() error
	Stop() error
}

type server struct {
	*http.Server
	started bool
	mutex   *sync.Mutex
}

// Get returns a Server
func Get(bindAddr string) Server {
	if s, ok := servers[bindAddr]; ok {
		return s
	}

	mutex.Lock()
	defer mutex.Unlock()

	if s, ok := servers[bindAddr]; ok {
		return s
	}

	servers[bindAddr] = &server{
		mutex: new(sync.Mutex),
		Server: &http.Server{
			Addr:         bindAddr,
			Handler:      pat.New(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	return servers[bindAddr]
}

// Router returns the gorilla/pat router
func (s *server) Router() *pat.Router {
	return s.Server.Handler.(*pat.Router)
}

func (s *server) Start() error {
	if s.started {
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	err := s.ListenAndServe()
	if err == nil {
		s.started = true
	}

	return err
}

func (s *server) Stop() error {
	// FIXME
	return errors.New("not implemented")
}

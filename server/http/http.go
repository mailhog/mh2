package http

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/pat"
	"github.com/ian-kent/service.go/handlers/requestID"
	"github.com/ian-kent/service.go/handlers/timeout"
	"github.com/ian-kent/service.go/log"
	"github.com/justinas/alice"
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
	router  *pat.Router
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

	router := pat.New()
	alice := alice.New(
		timeout.DefaultHandler,
		log.Handler,
		requestID.Handler(16),
	).Then(router)

	servers[bindAddr] = &server{
		mutex:  new(sync.Mutex),
		router: router,
		Server: &http.Server{
			Addr:         bindAddr,
			Handler:      alice,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	return servers[bindAddr]
}

// Router returns the gorilla/pat router
func (s *server) Router() *pat.Router {
	return s.router
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

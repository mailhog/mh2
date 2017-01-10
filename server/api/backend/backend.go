package backend

import (
	"github.com/mailhog/mh2/backend"

	// load backends
	_ "github.com/mailhog/mh2/backend/mongodb"
)

// API is the API backend
type API interface {
	// List returns a list of messages
	List(start, limit int) ([]*backend.Output, error)
	// Fetch returns a message based on message ID
	Fetch(backend.MessageID) (*backend.Output, error)
	// Close closes the backend
	Close() error
}

// New creates a new backend
func New(name string) (API, error) {
	return backend.New(name)
}

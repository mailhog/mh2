package backend

import (
	"fmt"
)

// MessageID represents a message ID
type MessageID string

// Backend represents a server backend
type Backend interface {
	// Receive receives the output of an SMTP conversation
	Receive(output *Output) error
	// Chan returns a channel to send notifications to
	Chan() chan MessageID
	// List returns a list of messages
	List(start, limit int) ([]*Output, error)
	// Fetch returns a message based on message ID
	Fetch(MessageID) (*Output, error)
	// Close closes the backend
	Close() error
}

var backends = map[string]func() (Backend, error){}

// Register registers a backend
func Register(name string, f func() (Backend, error)) {
	backends[name] = f
}

// New creates a new backend
func New(name string) (Backend, error) {
	if f, ok := backends[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("unknown backend: %s", name)
}

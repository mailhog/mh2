package backend

import (
	"github.com/mailhog/mh2/backend"

	// load backends
	_ "github.com/mailhog/mh2/backend/mongodb"
)

// SMTP is the SMTP backend
type SMTP interface {
	// Receive receives the output of an SMTP conversation
	Receive(output *backend.Output) error
	// Chan returns a channel to send notifications to
	Chan() chan backend.MessageID
	// Close closes the backend
	Close() error
}

// New creates a new backend
func New(name string) (SMTP, error) {
	return backend.New(name)
}

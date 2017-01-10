package backend

import (
	"fmt"

	"github.com/mailhog/mh2/backend/api"
	"github.com/mailhog/mh2/backend/smtp"
)

// Backend represents a server backend
type Backend interface {
	api.MessageReceiver
	api.MessageStorage
	smtp.OutputReceiver
	MessageReceiver() api.MessageReceiver
	MessageStorage() api.MessageStorage
	OutputReceiver() smtp.OutputReceiver
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

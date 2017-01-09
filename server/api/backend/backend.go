package backend

import "github.com/mailhog/mh2/backend/api"

// API is the API backend
type API interface {
	MessageReceiver() api.MessageReceiver
	MessageStorage() api.MessageStorage
	Close() error
}

package api

import "github.com/mailhog/mh2/backend/smtp"

// MessageID is a message ID
type MessageID string

// MessageStorage is a storage backend for SMTP messages
type MessageStorage interface {
	List(start, limit int) ([]*smtp.Output, error)
	Fetch(MessageID) (*smtp.Output, error)
}

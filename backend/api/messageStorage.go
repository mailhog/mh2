package api

import "github.com/mailhog/data"

// MessageID is a message ID
type MessageID string

// MessageStorage is a storage backend for SMTP messages
type MessageStorage interface {
	List(start, limit int) ([]*data.SMTPMessage, error)
	Fetch(MessageID) (*data.SMTPMessage, error)
}

package backend

import "github.com/mailhog/mh2/backend/smtp"

// SMTP is the SMTP backend
type SMTP interface {
	OutputReceiver() smtp.OutputReceiver
	Close() error
}

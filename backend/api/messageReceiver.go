package api

// MessageReceiver is a receiver of SMTP messages
type MessageReceiver interface {
	Chan() chan MessageID
}

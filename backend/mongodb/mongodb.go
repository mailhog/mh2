package mongodb

import (
	"errors"

	"github.com/mailhog/mh2/backend"
	"github.com/mailhog/mh2/backend/api"
	"github.com/mailhog/mh2/backend/smtp"
)

// b is the in-memory backend
type b struct {
	ch     chan api.MessageID
	exitCh chan int8
}

func init() {
	backend.Register("mongodb", New)
}

// New returns a new MongoDB backend
func New() (backend.Backend, error) {
	instance := &b{
		ch:     make(chan api.MessageID),
		exitCh: make(chan int8),
	}

	go func() {
		for {
			select {
			case <-instance.exitCh:
				break
			case _ = <-instance.ch:
				// TODO notification that message has arrived
			}
		}
	}()

	return instance, nil
}

// MessageReceiver implements api.MessageReceiver
func (b *b) MessageReceiver() api.MessageReceiver {
	return b
}

// MessageStorage implements api.MessageStorage
func (b *b) MessageStorage() api.MessageStorage {
	return b
}

func (b *b) OutputReceiver() smtp.OutputReceiver {
	return b
}

// Receive implements api.OutputReceiver
func (b *b) Receive(output *smtp.Output) error {
	return errors.New("not implemented")
}

// Close implements api.OutputReceiver and api.MessageReceiver
func (b *b) Close() error {
	b.exitCh <- 1
	return errors.New("not implemented")
}

// Chan implements api.MessageReceiver
func (b *b) Chan() chan api.MessageID {
	return b.ch
}

// List implements api.MessageStorage
func (b *b) List(start, limit int) ([]*smtp.Output, error) {
	return nil, errors.New("not implemented")
}

// Fetch implements api.MessageStorage
func (b *b) Fetch(api.MessageID) (*smtp.Output, error) {
	return nil, errors.New("not implemented")
}

package inmemory

import (
	"errors"
	"fmt"
	"sync"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/fatih/structs"
	"github.com/ian-kent/service.go/log"
	"github.com/mailhog/data"
	"github.com/mailhog/mh2/backend/api"
	"github.com/mailhog/mh2/backend/smtp"
)

var instance *b
var mtx sync.Mutex

// b is the in-memory backend
type b struct {
	ch       chan api.MessageID
	database *db.DB
	exitCh   chan int8
	messages *db.Col
}

// Backend is the in-memory backend
type Backend interface {
	api.MessageReceiver
	api.MessageStorage
	smtp.OutputReceiver
	MessageReceiver() api.MessageReceiver
	MessageStorage() api.MessageStorage
	OutputReceiver() smtp.OutputReceiver
	Close() error
}

// New returns a new in-memory backend
func New() (Backend, error) {
	mtx.Lock()
	defer mtx.Unlock()
	defer func() {
		log.Debug("backend", log.Data{"pointer": fmt.Sprintf("%p", instance)})
	}()

	if instance != nil {
		return instance, nil
	}

	_db, err := db.OpenDB("_db")
	if err != nil {
		return nil, err
	}

	instance = &b{
		database: _db,
		ch:       make(chan api.MessageID),
		exitCh:   make(chan int8),
	}

	instance.messages = _db.Use("messages")
	if instance.messages == nil {
		err = _db.Create("messages")
		if err != nil {
			return instance, err
		}
		instance.messages = _db.Use("messages")
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
	m := structs.Map(output)
	_, err := b.messages.Insert(m)
	return err
}

// Close implements api.OutputReceiver and api.MessageReceiver
func (b *b) Close() error {
	b.exitCh <- 1
	return b.database.Close()
}

// Chan implements api.MessageReceiver
func (b *b) Chan() chan api.MessageID {
	return b.ch
}

// List implements api.MessageStorage
func (b *b) List(start, limit int) ([]*data.SMTPMessage, error) {
	return nil, errors.New("not implemented")
}

// Fetch implements api.MessageStorage
func (b *b) Fetch(api.MessageID) (*data.SMTPMessage, error) {
	return nil, errors.New("not implemented")
}

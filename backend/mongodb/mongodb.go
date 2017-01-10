package mongodb

import (
	"errors"

	"github.com/ian-kent/service.go/log"
	"github.com/mailhog/mh2/backend"
	"gopkg.in/mgo.v2"
)

// b is the mongodb backend
type b struct {
	ch     chan backend.MessageID
	exitCh chan int8

	mongo        *mgo.Session
	messagesColl *mgo.Collection
}

func init() {
	backend.Register("mongodb", New)
}

// New returns a new mongodb backend
func New() (backend.Backend, error) {
	// FIXME make host/db/coll configurable
	host, db, coll := "localhost", "mh2", "conversations"

	mongoSession, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}
	log.Debug("connected to mongodb", log.Data{
		"host": host,
		"db":   db,
		"coll": coll,
	})

	instance := &b{
		ch:           make(chan backend.MessageID),
		exitCh:       make(chan int8),
		mongo:        mongoSession,
		messagesColl: mongoSession.DB(db).C(coll),
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

// Receive implements api.OutputReceiver
func (b *b) Receive(output *backend.Output) error {
	return b.messagesColl.Insert(output)
}

// Close implements api.OutputReceiver and api.MessageReceiver
func (b *b) Close() error {
	b.exitCh <- 1
	b.mongo.Close()
	return nil
}

// Chan implements api.MessageReceiver
func (b *b) Chan() chan backend.MessageID {
	return b.ch
}

// List implements api.MessageStorage
func (b *b) List(start, limit int) ([]*backend.Output, error) {
	return nil, errors.New("not implemented")
}

// Fetch implements api.MessageStorage
func (b *b) Fetch(backend.MessageID) (*backend.Output, error) {
	return nil, errors.New("not implemented")
}

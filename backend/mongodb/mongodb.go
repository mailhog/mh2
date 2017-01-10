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

	mongo             *mgo.Session
	conversationsColl *mgo.Collection
	messagesColl      *mgo.Collection
}

func init() {
	backend.Register("mongodb", New)
}

// New returns a new mongodb backend
func New() (backend.Backend, error) {
	// FIXME make host/db/cColl/mColl configurable
	host, db, cColl, mColl := "localhost", "mh2", "conversations", "messages"

	mongoSession, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}
	log.Debug("connected to mongodb", log.Data{
		"host":  host,
		"db":    db,
		"cColl": cColl,
		"mColl": mColl,
	})

	instance := &b{
		ch:                make(chan backend.MessageID),
		exitCh:            make(chan int8),
		mongo:             mongoSession,
		conversationsColl: mongoSession.DB(db).C(cColl),
		messagesColl:      mongoSession.DB(db).C(mColl),
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
	err := b.conversationsColl.Insert(output)
	if err != nil {
		return err
	}

	for _, m := range output.Messages {
		// TODO parse message and store properly
		// including a ref to the output context
		err = b.messagesColl.Insert(m)
		if err != nil {
			log.ErrorC(output.Context, err, nil)
		}
	}

	return nil
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

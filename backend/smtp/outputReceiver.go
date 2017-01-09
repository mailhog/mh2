package smtp

import (
	"sync"
	"time"

	"github.com/mailhog/data"
)

// OutputReceiver is a receiver of SMTP messages
type OutputReceiver interface {
	Receive(output *Output) error
}

type Output struct {
	Context       string
	RemoteAddress string
	TLS           bool
	Messages      []*data.SMTPMessage

	Data   []Data
	Proto  []Proto
	Events []Event

	LogData, LogProto, LogEvents bool

	mutex *sync.Mutex
}

func NewOutput(context string, remoteAddress string, logData, logEvents, logProto bool) *Output {
	return &Output{
		Context:       context,
		RemoteAddress: remoteAddress,
		LogData:       logData,
		LogEvents:     logEvents,
		LogProto:      logProto,

		mutex: new(sync.Mutex),
	}
}

type DataSender string

const Client DataSender = "CLIENT"
const Server DataSender = "SERVER"

type Data struct {
	Date   time.Time
	Sender DataSender
	Line   string
}

type Proto struct {
	Date  time.Time
	Event string
	Args  []interface{}
}

type Event struct {
	Date  time.Time
	Event string
	Args  []string
}

func (s *Output) RecordData(sender DataSender, line string) {
	if !s.LogData {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Data = append(s.Data, Data{time.Now(), sender, line})
}

func (s *Output) RecordProto(event string, args []interface{}) {
	if !s.LogProto {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Proto = append(s.Proto, Proto{time.Now(), event, args})
}

func (s *Output) RecordEvent(event string, args []string) {
	if !s.LogEvents {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Events = append(s.Events, Event{time.Now(), event, args})
}

func (s *Output) RecordMessage(message *data.SMTPMessage) {
	s.mutex.Lock()
	s.mutex.Unlock()
	s.Messages = append(s.Messages, message)
}

package backend

import (
	"sync"
	"time"

	"github.com/mailhog/data"
)

// Output represents an SMTP conversation
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

// NewOutput creates a new Output instance
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

// DataSender represents the data origin
type DataSender string

// Client represents data sent by the client
const Client DataSender = "CLIENT"

// Server represents data sent by the server
const Server DataSender = "SERVER"

// Data represents an individual line of the SMTP conversation
type Data struct {
	Date   time.Time
	Sender DataSender
	Line   string
}

// Proto represents data from the protocol state machine
type Proto struct {
	Date  time.Time
	Event string
	Args  []interface{}
}

// Event represents an event
type Event struct {
	Date  time.Time
	Event string
	Args  []string
}

// RecordData records a line of the SMTP conversation
func (s *Output) RecordData(sender DataSender, line string) {
	if !s.LogData {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Data = append(s.Data, Data{time.Now(), sender, line})
}

// RecordProto records data from the protocol state machine
func (s *Output) RecordProto(event string, args []interface{}) {
	if !s.LogProto {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Proto = append(s.Proto, Proto{time.Now(), event, args})
}

// RecordEvent records an event from the server
func (s *Output) RecordEvent(event string, args []string) {
	if !s.LogEvents {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Events = append(s.Events, Event{time.Now(), event, args})
}

// RecordMessage records a message received in the SMTP conversation
func (s *Output) RecordMessage(message *data.SMTPMessage) {
	s.mutex.Lock()
	s.mutex.Unlock()
	s.Messages = append(s.Messages, message)
}

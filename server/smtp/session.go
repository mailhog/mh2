package smtp

// http://www.rfc-editor.org/rfc/rfc5321.txt

import (
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/ian-kent/linkio"
	"github.com/ian-kent/service.go/log"
	"github.com/mailhog/data"
	"github.com/mailhog/smtp"
	"github.com/mailhog/storage"

	smtpbackend "github.com/mailhog/mh2/backend/smtp"
)

type acceptSession struct {
	remoteAddress string
	conn          io.ReadWriteCloser
	hostname      string
	monkey        ChaosMonkey

	messageChan chan *data.SMTPMessage
	outputChan  chan *smtpbackend.Output

	logProto, logData, logEvents bool
}

// Session represents a SMTP session using net.TCPConn
type Session struct {
	logData   bool
	sessionID string

	conn          io.ReadWriteCloser
	proto         *smtp.Protocol
	storage       storage.Storage
	remoteAddress string
	isTLS         bool
	line          string
	link          *linkio.Link

	messageChan chan *data.SMTPMessage
	outputChan  chan *smtpbackend.Output

	reader io.Reader
	writer io.Writer
	monkey ChaosMonkey

	output *smtpbackend.Output
}

// Accept starts a new SMTP session using io.ReadWriteCloser
func (s *smtpServer) Accept(a *acceptSession) {
	defer a.conn.Close()

	var context = ""
	if uuid, err := uuid.NewUUID(); err == nil {
		context = uuid.String()
	}

	proto := smtp.NewProtocol()
	proto.Hostname = a.hostname

	var link *linkio.Link
	reader := io.Reader(a.conn)
	writer := io.Writer(a.conn)

	if a.monkey != nil {
		linkSpeed := a.monkey.LinkSpeed(context)
		if linkSpeed != nil {
			link = linkio.NewLink(*linkSpeed * linkio.BytePerSecond)
			reader = link.NewLinkReader(io.Reader(a.conn))
			writer = link.NewLinkWriter(io.Writer(a.conn))
		}
	}

	session := &Session{
		sessionID:     context,
		conn:          a.conn,
		proto:         proto,
		messageChan:   a.messageChan,
		outputChan:    a.outputChan,
		remoteAddress: a.remoteAddress,
		isTLS:         false,
		line:          "",
		link:          link,
		reader:        reader,
		writer:        writer,
		monkey:        a.monkey,
		logData:       s.config.LogData,

		output: smtpbackend.NewOutput(context, a.remoteAddress, a.logData, a.logEvents, a.logProto),
	}

	proto.LogHandler = func(message string, args ...interface{}) {
		if s.config.LogProto {
			log.DebugC(context, "smtp protocol event", log.Data{"event": message, "args": args})
			session.output.RecordProto(message, args)
		}
	}
	proto.MessageReceivedHandler = session.acceptMessage
	proto.ValidateSenderHandler = session.validateSender
	proto.ValidateRecipientHandler = session.validateRecipient
	proto.ValidateAuthenticationHandler = session.validateAuthentication
	proto.GetAuthenticationMechanismsHandler = func() []string { return []string{"PLAIN"} }

	log.DebugC(context, "smtp: session started", nil)
	session.output.RecordEvent("session started", nil)
	session.Write(proto.Start())
	for session.Read() == true {
		if a.monkey != nil && a.monkey.Disconnect(context) {
			session.conn.Close()
			session.output.RecordEvent("monkey disconnected session", nil)
			break
		}
	}
	log.DebugC(context, "smtp: session finished", nil)
	session.output.RecordEvent("session finished", nil)
	if session.outputChan != nil {
		session.outputChan <- session.output
	}
}

func (c *Session) validateAuthentication(mechanism string, args ...string) (errorReply *smtp.Reply, ok bool) {
	c.output.RecordEvent("validating AUTH", append([]string{mechanism}, args...))
	if c.monkey != nil {
		ok := c.monkey.ValidAUTH(c.sessionID, mechanism, args...)
		if !ok {
			// FIXME better error?
			c.output.RecordEvent("monkey rejected AUTH", nil)
			return smtp.ReplyUnrecognisedCommand(), false
		}
	}
	c.output.RecordEvent("AUTH ok", nil)
	return nil, true
}

func (c *Session) validateRecipient(to string) bool {
	c.output.RecordEvent("validating RCPT TO", []string{to})
	if c.monkey != nil {
		ok := c.monkey.ValidRCPT(c.sessionID, to)
		if !ok {
			c.output.RecordEvent("monkey rejected RCPT TO", nil)
			return false
		}
	}
	c.output.RecordEvent("RCPT TO ok", nil)
	return true
}

func (c *Session) validateSender(from string) bool {
	c.output.RecordEvent("validating MAIL FROM", []string{from})
	if c.monkey != nil {
		ok := c.monkey.ValidMAIL(c.sessionID, from)
		if !ok {
			c.output.RecordEvent("monkey rejected MAIL FROM", nil)
			return false
		}
	}
	c.output.RecordEvent("MAIL FROM ok", nil)
	return true
}

func (c *Session) acceptMessage(msg *data.SMTPMessage) (id string, err error) {
	if c.logData {
		log.DebugC(c.sessionID, "smtp: accepting message", log.Data{"smtp_message": msg})
	} else {
		log.DebugC(c.sessionID, "smtp: accepting message", log.Data{"smtp_message": log.Data{
			"to":   msg.To,
			"from": msg.From,
			"helo": msg.Helo,
		}})
	}
	c.output.RecordMessage(msg)
	if c.messageChan != nil {
		c.messageChan <- msg
	}
	return
}

// Read reads from the underlying net.TCPConn
func (c *Session) Read() bool {
	buf := make([]byte, 1024)
	n, err := c.reader.Read(buf)

	if n == 0 {
		log.DebugC(c.sessionID, "smtp: connection closed by remote host", nil)
		io.Closer(c.conn).Close() // not sure this is necessary?
		c.output.RecordEvent("connection closed by remote host", nil)
		return false
	}

	if err != nil {
		log.ErrorC(c.sessionID, err, log.Data{"message": "smtp: error reading from socket"})
		c.output.RecordEvent("error reading from socket", nil)
		return false
	}

	text := string(buf[0:n])
	logText := strings.Replace(text, "\n", "\\n", -1)
	logText = strings.Replace(logText, "\r", "\\r", -1)
	if c.logData {
		log.DebugC(c.sessionID, "smtp: received data", log.Data{"length": n, "data": logText})
	}

	c.line += text

	for strings.Contains(c.line, "\r\n") {
		line, reply := c.proto.Parse(c.line)
		c.output.RecordData(smtpbackend.Client, c.line[:len(line)])
		c.line = line

		if reply != nil {
			c.Write(reply)
			if reply.Status == 221 {
				io.Closer(c.conn).Close()
				c.output.RecordEvent("connection closed following 221 reply", nil)
				return false
			}
		}
	}

	return true
}

// Write writes a reply to the underlying net.TCPConn
func (c *Session) Write(reply *smtp.Reply) {
	lines := reply.Lines()
	for _, l := range lines {
		logText := strings.Replace(l, "\n", "\\n", -1)
		logText = strings.Replace(logText, "\r", "\\r", -1)
		if c.logData {
			log.DebugC(c.sessionID, "smtp: sent data", log.Data{"length": len(l), "data": logText})
		}
		c.writer.Write([]byte(l))
		c.output.RecordData(smtpbackend.Server, l)
	}
}

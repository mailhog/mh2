package smtp

import (
	"io"
	"net"
	"sync"

	"github.com/ian-kent/gofigure"
	mh2backend "github.com/mailhog/mh2/backend"
	mh2server "github.com/mailhog/mh2/server"
	"github.com/mailhog/mh2/server/smtp/backend"

	"github.com/ian-kent/service.go/log"
)

type smtpServer struct {
	config   Config
	listener net.Listener
	exit     bool
	monkey   ChaosMonkey
	backend  backend.SMTP
}

// NewServer returns a new server
func NewServer() (mh2server.Server, error) {
	var smtpConfig = Config{
		BindAddr:            "0.0.0.0:1025",
		LogData:             false,
		LogProto:            false,
		Backend:             "mongodb",
		RecordSessionData:   true,
		RecordSessionEvents: true,
		RecordSessionProto:  true,
	}

	if err := gofigure.Gofigure(&smtpConfig); err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", smtpConfig.BindAddr)
	if err != nil {
		return nil, err
	}

	var jim ChaosMonkey
	if smtpConfig.Jim.Enabled {
		jim = &Jim{smtpConfig.Jim}
	}

	be, err := backend.New(smtpConfig.Backend)
	if err != nil {
		return nil, err
	}

	return &smtpServer{
		config:   smtpConfig,
		listener: listener,
		monkey:   jim,
		backend:  be,
	}, nil
}

// Start starts the server
func (s *smtpServer) Start() error {
	log.Debug("smtp: starting server", log.Data{"bind_addr": s.config.BindAddr})

	var wg sync.WaitGroup
	outputChan := make(chan *mh2backend.Output)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case output := <-outputChan:
				err := s.backend.Receive(output)
				if err != nil {
					if s.config.LogData {
						log.Error(err, log.Data{
							"output": output,
						})
					} else {
						log.Error(err, log.Data{
							"smtp_message": "<truncated>",
						})
					}
				}
			}
			if s.exit {
				break
			}
		}
	}()

	for {
		if s.exit {
			break
		}

		conn, err := s.listener.Accept()
		if err != nil {
			log.Error(err, log.Data{"message": "smtp: error accepting connection"})
			continue
		}

		go s.Accept(&acceptSession{
			remoteAddress: conn.(*net.TCPConn).RemoteAddr().String(),
			conn:          io.ReadWriteCloser(conn),
			outputChan:    outputChan,
			hostname:      s.config.Hostname,
			monkey:        s.monkey,
			logData:       s.config.RecordSessionData,
			logProto:      s.config.RecordSessionProto,
			logEvents:     s.config.RecordSessionEvents,
		})
	}

	wg.Wait()
	return nil
}

// Stop stops the server
func (s *smtpServer) Stop() error {
	s.exit = true
	return s.listener.Close()
}

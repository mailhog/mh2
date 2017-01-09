package server

import (
	"errors"
	"log"
	"sync"
)

// Server is a server
type Server interface {
	Start() error
	Stop() error
}

// New is a function which returns a new server
type New func() (Server, error)

// Start starts one or more servers
func Start(servers map[string]New) {
	errChan := make(chan error)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case err := <-errChan:
				log.Fatal(err)
			}
		}
	}()

	for k, v := range servers {
		srv, err := v()
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func(k string, srv Server) {
			defer wg.Done()
			err := srv.Start()
			if err != nil {
				errChan <- errors.New(k + ": " + err.Error())
			}
		}(k, srv)
	}

	wg.Wait()
}

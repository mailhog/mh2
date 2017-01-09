package main

import (
	"github.com/ian-kent/service.go/log"
	"github.com/mailhog/mh2/cmd"
	"github.com/mailhog/mh2/server"
	"github.com/mailhog/mh2/server/api"
	"github.com/mailhog/mh2/server/ui"
)

func main() {
	log.Namespace = "mh2-http"
	cmd.Main(func() {
		server.Start(map[string]server.New{
			"api": api.NewServer,
			"ui":  ui.NewServer,
		})
	})
}

package main

import (
	"github.com/ian-kent/service.go/log"
	"github.com/mailhog/mh2/cmd"
	"github.com/mailhog/mh2/server"
	"github.com/mailhog/mh2/server/smtp"
)

func main() {
	log.Namespace = "mh2-smtp"
	cmd.Main(func() {
		server.Start(map[string]server.New{
			"smtp": smtp.NewServer,
		})
	})
}

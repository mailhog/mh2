package inmemory

import (
	"github.com/mailhog/mh2/backend/api"
	"github.com/mailhog/mh2/backend/smtp"
	"github.com/mailhog/mh2/server/api/backend"
)

var _ backend.API = &b{}
var _ api.MessageReceiver = &b{}
var _ api.MessageStorage = &b{}
var _ smtp.OutputReceiver = &b{}

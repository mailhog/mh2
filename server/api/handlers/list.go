package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mailhog/mh2/server/api/backend"
)

// List is an API handler which lists messages
type List struct {
	API backend.API
}

// ServeHTTP implements http.Handler
func (l List) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	qStart := req.URL.Query().Get("start")
	qLimit := req.URL.Query().Get("limit")

	start, _ := strconv.Atoi(qStart)

	if start < 0 {
		start = 0
	}

	limit, _ := strconv.Atoi(qLimit)

	if limit > 1000 {
		limit = 1000
	} else if limit < 1 {
		limit = 1
	}

	messages, err := l.API.MessageStorage().List(start, limit)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := json.Marshal(&messages)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

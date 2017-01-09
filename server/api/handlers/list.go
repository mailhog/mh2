package handlers

import (
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

	start, err := strconv.Atoi(qStart)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	if start < 0 {
		start = 0
	}

	limit, err := strconv.Atoi(qLimit)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	if limit > 1000 {
		limit = 0
	}

	w.WriteHeader(200)
}

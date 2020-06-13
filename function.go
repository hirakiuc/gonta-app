package gonta

/**
 * This file provides the entry point for the Cloud Functions
 * (Use cmd/gonta/main.go for debug purpose)
 **/

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/server"
)

// nolint:gochecknoglobals
var srv *server.Gonta

// nolint:gochecknoinits
func init() {
	logger := log.GetLogger()
	srv = server.NewGonta(logger)
}

// Serve handles the http request.
func Serve(w http.ResponseWriter, r *http.Request) {
	srv.Serve(w, r)
}

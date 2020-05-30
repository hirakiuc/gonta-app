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

// Serve handles the http request.
func Serve(w http.ResponseWriter, r *http.Request) {
	log := log.GetLogger()

	gonta := server.NewGonta(log)

	gonta.Serve(w, r)
}

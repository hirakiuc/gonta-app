package main

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/event"
	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/server"

	"go.uber.org/zap"
)

func main() {
	logger := log.GetLogger()

	d := event.NewDispatcher(logger)
	srv := server.NewGonta(logger, d)

	http.HandleFunc("/serve", srv.Serve)

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		logger.Fatal("Failed", zap.Error(err))
	}
}

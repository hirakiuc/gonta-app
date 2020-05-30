package main

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/server"
	"go.uber.org/zap"
)

func main() {
	logger := log.GetLogger()

	srv := server.NewGonta(logger)

	http.HandleFunc("/serve", srv.Serve)

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		logger.Fatal("Failed", zap.Error(err))
	}
}

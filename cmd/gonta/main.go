package main

import (
	"net/http"
	"os"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/event"
	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/server"

	"go.uber.org/zap"
)

func main() {
	logger := log.GetLogger()

	conf := config.NewConfig()

	err := conf.Load()
	if err != nil {
		logger.Error("Failed to load config", zap.Error(err))
		os.Exit(1)
	}

	d := event.NewDispatcher(logger, conf)
	srv := server.NewGonta(logger, d, conf)

	http.HandleFunc("/serve", srv.SlackVerify(srv.ServeEvents))

	http.HandleFunc("/events", srv.SlackVerify(srv.ServeEvents))
	http.HandleFunc("/actions", srv.SlackVerify(srv.ServeActions))
	http.HandleFunc("/commands", srv.SlackVerify(srv.ServeCommands))

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		logger.Fatal("Failed", zap.Error(err))
	}
}

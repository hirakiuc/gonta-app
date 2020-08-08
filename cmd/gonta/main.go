package main

import (
	"net/http"
	"os"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/event/data"
	"github.com/hirakiuc/gonta-app/event/queue"
	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/server"
	"github.com/hirakiuc/gonta-app/usecase"

	"go.uber.org/zap"
)

const QueueSize = 50

func main() {
	logger := log.GetLogger()

	conf := config.NewConfig()

	err := conf.Load()
	if err != nil {
		logger.Error("Failed to load config", zap.Error(err))
		os.Exit(1)
	}

	q := queue.New(QueueSize, logger)
	d := data.NewProvider()
	usecase.Configure(q, d, conf.HandlerConfig(), logger)

	go q.Start()

	srv := server.NewGonta(logger, conf, q, d)
	defer srv.Wait()

	http.HandleFunc("/serve", srv.SlackVerify(srv.ServeEvents))

	http.HandleFunc("/health", srv.ServeHealth)
	http.HandleFunc("/events", srv.SlackVerify(srv.ServeEvents))
	http.HandleFunc("/actions", srv.SlackVerify(srv.ServeActions))
	http.HandleFunc("/commands", srv.SlackVerify(srv.ServeCommands))
	http.HandleFunc("/data", srv.SlackVerify(srv.ServeData))

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		logger.Fatal("Failed", zap.Error(err))
	}
}

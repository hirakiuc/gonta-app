package main

import (
	"net/http"
	"os"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/queue"
	"github.com/hirakiuc/gonta-app/server"
	"github.com/hirakiuc/gonta-app/usecase"

	"go.uber.org/zap"
)

const QueueSize = 50

func configureCallbacks(q *queue.Queue, conf *config.Config, logger *zap.Logger) {
	echo := usecase.NewEcho(conf, logger)
	q.AddEventCallback(echo.ReceiveEvent)
}

func main() {
	logger := log.GetLogger()

	conf := config.NewConfig()

	err := conf.Load()
	if err != nil {
		logger.Error("Failed to load config", zap.Error(err))
		os.Exit(1)
	}

	q := queue.New(QueueSize, logger)
	configureCallbacks(q, conf, logger)

	go q.Start()

	defer func() {
		q.Stop()
		q.WaitUntilFinish()
	}()

	srv := server.NewGonta(logger, conf, q)

	http.HandleFunc("/serve", srv.SlackVerify(srv.ServeEvents))

	http.HandleFunc("/health", srv.ServeHealth)
	http.HandleFunc("/events", srv.SlackVerify(srv.ServeEvents))
	http.HandleFunc("/actions", srv.SlackVerify(srv.ServeActions))
	http.HandleFunc("/commands", srv.SlackVerify(srv.ServeCommands))

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		logger.Fatal("Failed", zap.Error(err))
	}
}

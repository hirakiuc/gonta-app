package main

import (
	"net/http"
	"os"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/event/data"
	"github.com/hirakiuc/gonta-app/event/queue"
	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/server"
	"github.com/hirakiuc/gonta-app/usecase/release"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

const QueueSize = 50

func configure(q *queue.Queue, d *data.Provider, conf *config.HandlerConfig, logger *zap.Logger) {
	// Configure release callbacks
	rel := release.New(conf, logger)

	q.AddEventCallback(slackevents.AppMention, rel.Start)

	d.AddProvider(release.SelectVersionBlockID, rel.FetchVersions)

	q.AddBlockActionCallback(
		slack.InteractionTypeBlockActions,
		release.SelectVersionBlockID,
		rel.ConfirmRelease,
	)

	q.AddBlockActionCallback(
		slack.InteractionTypeBlockActions,
		release.ConfirmDeploymentBlockID,
		rel.InvokeRelease,
	)
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
	d := data.NewProvider()
	configure(q, d, conf.HandlerConfig(), logger)

	go q.Start()

	srv := server.NewGonta(logger, conf, q, d)
	defer srv.Wait()

	http.HandleFunc("/health", srv.ServeHealth)

	// Endpoint for Event Subscriptions
	http.HandleFunc("/events", srv.SlackVerify(srv.ServeEvents))

	// Endpoint for Interactivity
	http.HandleFunc("/actions", srv.SlackVerify(srv.ServeActions))

	// Endpoint for the slash Commands
	http.HandleFunc("/commands", srv.SlackVerify(srv.ServeCommands))

	// Endpoint for the external data feature of the select menus.
	http.HandleFunc("/data", srv.SlackVerify(srv.ServeData))

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		logger.Fatal("Failed", zap.Error(err))
	}
}

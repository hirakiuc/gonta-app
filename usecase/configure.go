package usecase

import (
	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/queue"

	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

func Configure(q *queue.Queue, conf *config.HandlerConfig, logger *zap.Logger) {
	echo := NewEcho(conf, logger)
	q.AddEventCallback(slackevents.AppMention, echo.ReceiveEvent)
}

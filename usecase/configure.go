package usecase

import (
	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/queue"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

func Configure(q *queue.Queue, conf *config.HandlerConfig, logger *zap.Logger) {
	// Configure echo callbacks
	echo := NewEcho(conf, logger)
	q.AddEventCallback(slackevents.AppMention, echo.ReceiveEvent)

	// Configure release callbacks
	rel := NewRelease(conf, logger)

	q.AddEventCallback(slackevents.AppMention, rel.ShowVersionChooser)

	q.AddBlockActionCallback(
		slack.InteractionTypeBlockActions,
		selectVersionActionID,
		rel.ConfirmRelease,
	)

	q.AddBlockActionCallback(
		slack.InteractionTypeBlockActions,
		confirmDeploymentActionID,
		rel.InvokeRelease,
	)
}

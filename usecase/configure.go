package usecase

import (
	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/event/data"
	"github.com/hirakiuc/gonta-app/event/queue"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

func Configure(q *queue.Queue, d *data.Provider, conf *config.HandlerConfig, logger *zap.Logger) {
	// Configure release callbacks
	rel := NewRelease(conf, logger)

	d.AddProvider(selectVersionBlockID, rel.FetchVersions)

	q.AddEventCallback(slackevents.AppMention, rel.ShowVersionChooser)

	q.AddBlockActionCallback(
		slack.InteractionTypeBlockActions,
		selectVersionBlockID,
		rel.ConfirmRelease,
	)

	q.AddBlockActionCallback(
		slack.InteractionTypeBlockActions,
		confirmDeploymentBlockID,
		rel.InvokeRelease,
	)
}

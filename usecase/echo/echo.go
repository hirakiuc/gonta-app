package usecase

import (
	"context"
	"fmt"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/usecase"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

type Echo struct {
	usecase.Base
}

func New(c *config.HandlerConfig, logger *zap.Logger) *Echo {
	return &Echo{
		Base: usecase.Base{
			Config: c,
			Logger: logger,
		},
	}
}

func (u *Echo) ReceiveEvent(e *slackevents.EventsAPIEvent) error {
	u.Logger.Info("Receive event:handler-echo", zap.String("handler", "echo"))

	ev, err := u.CastAppMentionEvent(e)
	if err != nil {
		return err
	}

	api := u.SlackAPI()

	ctx := context.Background()
	msg := slack.MsgOptionText("Yes, hello", false)

	u.Logger.Info("Start sending a message...")

	channel, tstamp, err := api.PostMessageContext(ctx, ev.Channel, msg)
	if err != nil {
		msg := fmt.Sprintf("Failed to send a message:%s", err.Error())
		u.Logger.Error(msg, zap.Error(err))

		return err
	}

	u.Logger.Debug("Sent a message", zap.String("channel", channel), zap.String("timestamp", tstamp))

	return nil
}

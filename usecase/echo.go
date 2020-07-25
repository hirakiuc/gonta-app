package usecase

import (
	"context"
	"fmt"

	"github.com/hirakiuc/gonta-app/config"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

type Echo struct {
	Base
}

func NewEcho(c *config.HandlerConfig, logger *zap.Logger) *Echo {
	return &Echo{
		Base: Base{
			config: c,
			logger: logger,
		},
	}
}

func (u *Echo) ReceiveEvent(e *slackevents.EventsAPIEvent) error {
	u.logger.Info("Receive event:handler-echo", zap.String("handler", "echo"))

	ev, err := castAppMentionEvent(e)
	if err != nil {
		return err
	}

	api := u.slackAPI()

	ctx := context.Background()
	msg := slack.MsgOptionText("Yes, hello", false)

	u.logger.Info("Start sending a message...")

	channel, tstamp, err := api.PostMessageContext(ctx, ev.Channel, msg)
	if err != nil {
		msg := fmt.Sprintf("Failed to send a message:%s", err.Error())
		u.logger.Error(msg, zap.Error(err))

		return err
	}

	u.logger.Debug("Sent a message", zap.String("channel", channel), zap.String("timestamp", tstamp))

	return nil
}

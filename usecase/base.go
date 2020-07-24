package usecase

import (
	"github.com/hirakiuc/gonta-app/config"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type Base struct {
	config *config.HandlerConfig
	logger *zap.Logger
}

func (b *Base) slackAPI() *slack.Client {
	return slack.New(b.config.BotAccessToken)
}

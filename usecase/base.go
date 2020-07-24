package usecase

import (
	"crypto/tls"
	"net/http"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type Base struct {
	config *config.HandlerConfig
	logger *zap.Logger
}

func (b *Base) slackAPI() *slack.Client {
	// nolint:gosec
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return slack.New(b.config.BotAccessToken, slack.OptionHTTPClient(c))
}

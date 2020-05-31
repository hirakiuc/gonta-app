package handler

import (
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

// ClientOption describe a config for slack.Client.
type ClientOption struct {
	Token string
	Log   *zap.Logger
}

// GetClient return a slack.Client instance.
func GetClient(opts ClientOption) (*slack.Client, error) {
	return slack.New(opts.Token), nil
}

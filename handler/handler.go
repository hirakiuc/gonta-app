package handler

import (
	"context"

	"github.com/hirakiuc/gonta-app/config"

	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

type Handler interface {
	SetLogger(logger *zap.Logger)
	SetConfig(c *config.HandlerConfig)
	Handle(ctx context.Context, event *slackevents.EventsAPIEvent) error
}

type BaseHandler struct {
	log    *zap.Logger
	config *config.HandlerConfig
}

func (h *BaseHandler) SetLogger(logger *zap.Logger) {
	h.log = logger
}

func (h *BaseHandler) SetConfig(c *config.HandlerConfig) {
	h.config = c
}

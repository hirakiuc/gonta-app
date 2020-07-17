package server

import (
	"context"
	"net/http"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

type Handler interface {
	SetLogger(logger *zap.Logger)
	SetConfig(c *config.Config)
	Handle(ctx context.Context, w http.ResponseWriter, event *slackevents.EventsAPIEvent) error
}

type BaseHandler struct {
	log    *zap.Logger
	config *config.Config
}

func (h *BaseHandler) SetLogger(logger *zap.Logger) {
	h.log = logger
}

func (h *BaseHandler) SetConfig(c *config.Config) {
	h.config = c
}

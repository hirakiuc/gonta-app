package handler

import (
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

type Handler interface {
	SetLogger(logger *zap.Logger)
	Handle(event *slackevents.EventsAPIEvent) error
}

type BaseHandler struct {
	log *zap.Logger
}

func (h *BaseHandler) SetLogger(logger *zap.Logger) {
	h.log = logger
}

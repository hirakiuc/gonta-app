package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

var ErrUnsupportedEventType = errors.New("unsuppoted event type")

type CallbackEventHandler struct {
	BaseHandler
}

func NewCallbackEventHandler() *CallbackEventHandler {
	return &CallbackEventHandler{}
}

func (h *CallbackEventHandler) Handle(w http.ResponseWriter, event *slackevents.EventsAPIEvent) error {
	log := h.log

	innerEvent := event.InnerEvent
	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		handler := NewMentionHandler()
		handler.SetLogger(log)

		return handler.Handle(w, event, ev)
	default:
		log.Error("un-supported event type", zap.String("type", innerEvent.Type))

		return fmt.Errorf(
			"unsupported event type:%s %w",
			innerEvent.Type,
			ErrUnsupportedEventType,
		)
	}
}

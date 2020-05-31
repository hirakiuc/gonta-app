package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hirakiuc/gonta-app/model"
	"go.uber.org/zap"
)

var ErrUnsupportedEventType = errors.New("unsuppoted event type")

type CallbackEventHandler struct {
	BaseHandler
}

func NewCallbackEventHandler() *CallbackEventHandler {
	return &CallbackEventHandler{}
}

func (h *CallbackEventHandler) Handle(w http.ResponseWriter, msg *model.CallbackEvent) error {
	log := h.log

	t, err := msg.GetEventType()
	if err != nil {
		log.Error("Failed to extract event type", zap.Error(err))

		return err
	}

	switch t {
	case "app_mention":
		handler := NewMentionHandler()
		handler.SetLogger(log)

		return handler.Handle(w, msg)
	default:
		log.Error("un-supported event type", zap.String("type", t))

		return fmt.Errorf(
			"unsupported event type:%s %w",
			t,
			ErrUnsupportedEventType,
		)
	}
}

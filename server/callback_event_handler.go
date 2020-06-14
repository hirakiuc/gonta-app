package server

import (
	"errors"
	"net/http"

	"github.com/hirakiuc/gonta-app/event"

	"github.com/slack-go/slack/slackevents"
)

var ErrUnsupportedEventType = errors.New("unsuppoted event type")

type CallbackEventHandler struct {
	dispatcher *event.Dispatcher

	BaseHandler
}

func NewCallbackEventHandler(dispatcher *event.Dispatcher) *CallbackEventHandler {
	return &CallbackEventHandler{
		dispatcher: dispatcher,
	}
}

func (h *CallbackEventHandler) Handle(w http.ResponseWriter, e *slackevents.EventsAPIEvent) error {
	wg := (h.dispatcher).Dispatch(e)
	wg.Wait()

	w.WriteHeader(http.StatusOK)

	return nil
}

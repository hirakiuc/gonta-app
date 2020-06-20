package server

import (
	"context"
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

func (h *CallbackEventHandler) Handle(ctx context.Context, w http.ResponseWriter, e *slackevents.EventsAPIEvent) error {
	// Dispatch this event to each registered handlers.
	wg := (h.dispatcher).Dispatch(ctx, e)

	// Wait until invoked handlers finish
	wg.Wait()

	// Send a response (200 OK)
	w.WriteHeader(http.StatusOK)

	return nil
}

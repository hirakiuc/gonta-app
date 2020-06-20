package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

var ErrUnexpectedEventType = errors.New("unexpected event type")

// MentionHandler describe a instance of MentionHandler.
type MentionHandler struct {
	BaseHandler
}

// NewMentionHandler return an instance of MentionHandler.
func NewMentionHandler() *MentionHandler {
	return &MentionHandler{}
}

func (h *MentionHandler) innerEvent(e *slackevents.EventsAPIEvent) (*slackevents.AppMentionEvent, error) {
	innerEvent := e.InnerEvent

	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		return ev, nil
	default:
		err := fmt.Errorf("%w", ErrUnexpectedEventType)
		h.log.Error("unexpected event type", zap.Error(err))

		return nil, err
	}
}

// Reply send a response to the slack.
func (h *MentionHandler) Handle(ctx context.Context, e *slackevents.EventsAPIEvent) error {
	log := h.log

	ev, err := h.innerEvent(e)
	if err != nil {
		return err
	}

	log.Debug("MentionHandler handle", zap.String("received", ev.Text))

	/*
		api := slack.New(e.Token)

		channel, tstamp, err := api.PostMessageContext(ctx, ev.Channel, slack.MsgOptionText("Yes, hello", false))
		if err != nil {
			log.Error("failed to send a message", zap.Error(err))
			return err
		}

		log.Debug("sent a message", zap.String("channel", channel), zap.String("timestamp", tstamp))
	*/

	return nil
}

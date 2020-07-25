package usecase

import (
	"errors"
	"fmt"

	"github.com/slack-go/slack/slackevents"
)

var ErrUnexpectedEventType = errors.New("unexpected event type")

func castAppMentionEvent(e *slackevents.EventsAPIEvent) (*slackevents.AppMentionEvent, error) {
	innerEvent := e.InnerEvent

	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		return ev, nil
	default:
		return nil, fmt.Errorf("%w", ErrUnexpectedEventType)
	}
}

package handler

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

const (
	actionSelect = "select"
	// actionStart  = "start"
	// actionCancel = "cancel"
)

// BeerSelectHandler describe a instance of BeerSelect Request.
type BeerSelectHandler struct {
	BaseHandler
}

// NewBeerSelectHandler return an BeerSelectReply instance.
func NewBeerSelectHandler() *BeerSelectHandler {
	return &BeerSelectHandler{}
}

func (h *BeerSelectHandler) innerEvent(e *slackevents.EventsAPIEvent) (*slackevents.AppMentionEvent, error) {
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

// Handler a beer select event.
func (h *BeerSelectHandler) Handle(_ context.Context, e *slackevents.EventsAPIEvent) error {
	log := h.log

	log.Debug("BaseSelectReplyer reply:empty")

	ev, err := h.innerEvent(e)
	if err != nil {
		return err
	}

	api := slackAPI()

	attachment := slack.Attachment{
		Text:       "Which beer do you want? :beer:",
		Color:      "#f9a41b",
		CallbackID: "beer",
		Actions: []slack.AttachmentAction{
			{
				Name: actionSelect,
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "Asahi Super Dry",
						Value: "Asahi Super Dry",
					},
					{
						Text:  "Kirin Lager Beer",
						Value: "Kirin Lager Beer",
					},
					{
						Text:  "Sapporo Black Label",
						Value: "Sapporo Black Label",
					},
					{
						Text:  "Suntory Malts",
						Value: "Suntory Malts",
					},
					{
						Text:  "Yona Yona Ale",
						Value: "Yona Yona Ale",
					},
				},
			},
		},
	}

	channelID, tstamp, err := api.PostMessage(
		ev.Channel,
		slack.MsgOptionText("Some Text", false),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		log.Error("BaseSelectReplyer failed", zap.Error(err))
	}

	log.Debug("Post Message", zap.String("channelID", channelID), zap.String("timestamp", tstamp))

	return nil
}

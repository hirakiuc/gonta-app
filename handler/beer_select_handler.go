package handler

import (
	"context"

	"github.com/slack-go/slack/slackevents"
)

/*
const (
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)
*/

// BeerSelectHandler describe a instance of BeerSelect Request.
type BeerSelectHandler struct {
	BaseHandler
}

// NewBeerSelectHandler return an BeerSelectReply instance.
func NewBeerSelectHandler() *BeerSelectHandler {
	return &BeerSelectHandler{}
}

// Handler a beer select event.
func (h *BeerSelectHandler) Handle(_ context.Context, event *slackevents.EventsAPIEvent) error {
	log := h.log

	log.Debug("BaseSelectReplyer reply:empty")
	/*
		client, err := GetClient()
		if err != nil {
			// Can't respond to the event.
			w.WriteHeader(http.StatusOK)
			return
		}

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

			channelID, timestamp, err := client.PostMessage(
				e.ChannelID,
				slack.MsgOptionText("Some Text", false),
				slack.MsgOptionAttachments(attachment),
			)
			if err != nil {
				log.Error("BaseSelectReplyer failed", zap.Error(err))
			}

			log.Debug("Post Message", zap.String("channelID", channelID), zap.String("timestamp", timestamp))
	*/

	return nil
}

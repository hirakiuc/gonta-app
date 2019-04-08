package reply

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/event"
	"github.com/labstack/gommon/log"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

const (
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

// BeerSelectReply describe a instance of BeerSelect Request.
type BeerSelectReplyer struct{}

// NewBeerSelectReply return an BeerSelectReply instance.
func NewBeerSelectReplyer() *BeerSelectReplyer {
	return &BeerSelectReplyer{}
}

// Reply send a beer select reply to the slack channel
func (req *BeerSelectReplyer) Reply(w http.ResponseWriter, msg *event.CallbackEvent) {
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

	channelID, timestamp, err := client.PostMessage(msg.ChannelID, slack.MsgOptionText("Some Text", false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		log.Error("BaseSelectReplyer failed", zap.Error(err))
	}

	log.Debug("Post Message", zap.String("channelID", channelID), zap.String("timestamp", timestamp))
}

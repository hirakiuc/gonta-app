package handler

import (
	"github.com/slack-go/slack/slackevents"
)

func GenerateHandlerMap() map[string][]Handler {
	return map[string][]Handler{
		slackevents.AppMention: {
			NewMentionHandler(),
		},
	}
}

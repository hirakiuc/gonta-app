package handler

import (
	"github.com/slack-go/slack/slackevents"
)

// MentionHandler describe a instance of MentionHandler.
type MentionHandler struct {
	BaseHandler
}

// NewMentionHandler return an instance of MentionHandler.
func NewMentionHandler() *MentionHandler {
	return &MentionHandler{}
}

// Reply send a response to the slack.
func (h *MentionHandler) Handle(e *slackevents.EventsAPIEvent) error {
	log := h.log

	log.Debug("MentionHandler handle")

	return nil
}

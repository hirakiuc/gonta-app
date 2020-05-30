package handler

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/event"
	"github.com/hirakiuc/gonta-app/log"
)

// MentionHandler describe a instance of MentionHandler.
type MentionHandler struct{}

// NewMentionHandler return an instance of MentionHandler.
func NewMentionHandler() *MentionHandler {
	return &MentionHandler{}
}

// Reply send a response to the slack.
func (req *MentionHandler) Handle(w http.ResponseWriter, msg *event.CallbackEvent) error {
	logger := log.GetLogger()
	logger.Debug("MentionHandler handle")
	w.WriteHeader(http.StatusOK)

	return nil
}

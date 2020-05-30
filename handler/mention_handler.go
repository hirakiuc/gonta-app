package handler

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/model"
)

// MentionHandler describe a instance of MentionHandler.
type MentionHandler struct{}

// NewMentionHandler return an instance of MentionHandler.
func NewMentionHandler() *MentionHandler {
	return &MentionHandler{}
}

// Reply send a response to the slack.
func (req *MentionHandler) Handle(w http.ResponseWriter, msg *model.CallbackEvent) error {
	logger := log.GetLogger()
	logger.Debug("MentionHandler handle")
	w.WriteHeader(http.StatusOK)

	return nil
}

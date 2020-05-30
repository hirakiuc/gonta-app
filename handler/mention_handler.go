package handler

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/model"
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
func (h *MentionHandler) Handle(w http.ResponseWriter, msg *model.CallbackEvent) error {
	log := h.log

	log.Debug("MentionHandler handle")
	w.WriteHeader(http.StatusOK)

	return nil
}

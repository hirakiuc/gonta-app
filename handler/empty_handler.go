package handler

import (
	"net/http"

	event "github.com/hirakiuc/gonta-app/event"
)

// EmptyHandler describe a replyer with empty response.
type EmptyHandler struct{}

// NewEmptyHandler return an instance of EmptyHandler.
func NewEmptyHandler() *EmptyHandler {
	return &EmptyHandler{}
}

// Reply respond an empty response.
func (r *EmptyHandler) Reply(w http.ResponseWriter, msg *event.CallbackEvent) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

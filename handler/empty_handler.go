package handler

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/model"
)

// EmptyHandler describe a replyer with empty response.
type EmptyHandler struct {
	BaseHandler
}

// NewEmptyHandler return an instance of EmptyHandler.
func NewEmptyHandler() *EmptyHandler {
	return &EmptyHandler{}
}

// Reply respond an empty response.
func (h *EmptyHandler) Handle(w http.ResponseWriter, msg *model.CallbackEvent) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

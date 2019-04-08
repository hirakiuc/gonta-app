package reply

import (
	"net/http"

	event "github.com/hirakiuc/gonta-app/event"
)

// EmptyReplyer describe a replyer with empty response.
type EmptyReplyer struct{}

// NewEmptyReplyer return an instance of EmptyReplyer.
func NewEmptyReplyer() *EmptyReplyer {
	return &EmptyReplyer{}
}

// Reply respond an empty response.
func (r *EmptyReplyer) Reply(w http.ResponseWriter, msg *event.CallbackEvent) {
	w.WriteHeader(http.StatusOK)
}

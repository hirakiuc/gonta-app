package reply

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/event"
	"github.com/hirakiuc/gonta-app/log"
)

// MentionReplyer describe a instance of MentionReplyer.
type MentionReplyer struct{}

// NewMentionReplyer return an instance of MentionReplyer.
func NewMentionReplyer() *MentionReplyer {
	return &MentionReplyer{}
}

// Reply send a response to the slack.
func (req *MentionReplyer) Reply(w http.ResponseWriter, msg *event.CallbackEvent) {
	logger := log.GetLogger()
	logger.Debug("MentionReplyer reply")
	w.WriteHeader(http.StatusOK)
}

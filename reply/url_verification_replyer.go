package reply

import (
	"encoding/json"
	"net/http"

	"github.com/hirakiuc/gonta-app/event"
	"github.com/hirakiuc/gonta-app/log"
	"go.uber.org/zap"
)

type challengeResponse struct {
	Challenge string `json:"challenge"`
}

// URLVerificationReplyer describe a instance of URLVerification replyer.
type URLVerificationReplyer struct{}

// NewURLVerificationReplyer return an URLVerificationReply instance.
func NewURLVerificationReplyer() *URLVerificationReplyer {
	return &URLVerificationReplyer{}
}

// Reply send the response for the URLVerification reply.
func (replyer *URLVerificationReplyer) Reply(w http.ResponseWriter, msg *event.URLVerificationEvent) {
	log := log.GetLogger()

	challenge := challengeResponse{Challenge: msg.Challenge}

	res, err := json.Marshal(challenge)
	if err != nil {
		log.Error("Failed to marshal json response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

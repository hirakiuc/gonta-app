package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/hirakiuc/gonta-app/event"
	"github.com/hirakiuc/gonta-app/log"
	"go.uber.org/zap"
)

type challengeResponse struct {
	Challenge string `json:"challenge"`
}

// URLVerificationHandler describe a instance of URLVerification replyer.
type URLVerificationHandler struct{}

// NewURLVerificationHandler return an URLVerificationReply instance.
func NewURLVerificationHandler() *URLVerificationHandler {
	return &URLVerificationHandler{}
}

// Reply send the response for the URLVerification reply.
func (replyer *URLVerificationHandler) Handle(w http.ResponseWriter, msg *event.URLVerificationEvent) {
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

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		// Just log error
		fmt.Fprintf(os.Stderr, "Failed: %s", err)
	}
}

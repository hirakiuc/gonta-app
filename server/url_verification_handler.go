package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

var ErrUnexpectedEvent = errors.New("unexpected event type")

// URLVerificationHandler describe a instance of URLVerification replyer.
type URLVerificationHandler struct {
	BaseHandler
}

// NewURLVerificationHandler return an URLVerificationReply instance.
func NewURLVerificationHandler() *URLVerificationHandler {
	return &URLVerificationHandler{}
}

// Reply send the response for the URLVerification reply.
func (h *URLVerificationHandler) Handle(w http.ResponseWriter, event *slackevents.EventsAPIEvent) error {
	log := h.log

	d, ok := event.Data.(slackevents.EventsAPIURLVerificationEvent)
	if !ok {
		log.Error("Unexpected type")
		w.WriteHeader(http.StatusInternalServerError)

		return fmt.Errorf("unexpected event type:%w", ErrUnexpectedEvent)
	}

	challenge := slackevents.ChallengeResponse{
		Challenge: d.Challenge,
	}

	res, err := json.Marshal(challenge)
	if err != nil {
		log.Error("Failed to marshal json response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error("Failed to respond the result")
		w.WriteHeader(http.StatusInternalServerError)

		return err
	}

	return nil
}

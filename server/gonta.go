package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hirakiuc/gonta-app/event"

	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

var ErrUnexpectedEventType = errors.New("unexpected event type")

// Gonta describe a http server to serve gonta services.
type Gonta struct {
	log        *zap.Logger
	dispatcher *event.Dispatcher
}

func NewGonta(logger *zap.Logger, d *event.Dispatcher) *Gonta {
	return &Gonta{
		log:        logger,
		dispatcher: d,
	}
}

// Serve handles the http request.
func (s *Gonta) Serve(w http.ResponseWriter, r *http.Request) {
	log := s.log

	if r.Method != http.MethodPost {
		log.Debug("Invalid http method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	opts := slackevents.OptionVerifyToken(
		&slackevents.TokenComparator{
			VerificationToken: getVerificationToken(),
		},
	)

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), opts)
	if err != nil {
		log.Error("Failed to parse request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	handler, err := s.handlerByEventType(eventsAPIEvent.Type)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	ctx := context.Background()

	handler.SetLogger(log)

	if err = handler.Handle(ctx, w, &eventsAPIEvent); err != nil {
		log.Error("Failed to process the request", zap.Error(err))
	}
}

func getVerificationToken() string {
	return os.Getenv("VERIFICATION_TOKEN")
}

func (s *Gonta) handlerByEventType(eventType string) (Handler, error) {
	log := s.log

	switch eventType {
	case slackevents.URLVerification:
		return NewURLVerificationHandler(), nil
	case slackevents.CallbackEvent:
		return NewCallbackEventHandler(s.dispatcher), nil
	default:
		log.Error("Unexpected event type", zap.String("type", eventType))

		return nil, fmt.Errorf("unexpected event type:%s %w", eventType, ErrUnexpectedEventType)
	}
}

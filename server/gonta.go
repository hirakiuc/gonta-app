package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/event"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

var ErrUnexpectedEventType = errors.New("unexpected event type")

// Gonta describe a http server to serve gonta services.
type Gonta struct {
	log        *zap.Logger
	dispatcher *event.Dispatcher
	config     *config.Config
}

func NewGonta(logger *zap.Logger, d *event.Dispatcher, c *config.Config) *Gonta {
	return &Gonta{
		log:        logger,
		dispatcher: d,
		config:     c,
	}
}

// nolint:interfacer
func (s *Gonta) SlackVerify(next http.HandlerFunc) http.HandlerFunc {
	log := s.log

	return func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_SIGNING_SECRET"))
		if err != nil {
			log.Error("Failed to create verifier", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		bodyReader := io.TeeReader(r.Body, &verifier)

		body, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			log.Error("Failed to read body", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = verifier.Ensure()
		if err != nil {
			log.Error("Failed to verify request", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		next.ServeHTTP(w, r)
	}
}

// ServeEvents handles the http request.
func (s *Gonta) ServeEvents(w http.ResponseWriter, r *http.Request) {
	log := s.log

	if r.Method != http.MethodPost {
		log.Debug("Invalid http method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	opts := slackevents.OptionNoVerifyToken()

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
	handler.SetConfig(s.config)

	if err = handler.Handle(ctx, w, &eventsAPIEvent); err != nil {
		log.Error("Failed to process the request", zap.Error(err))
	}
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

func (s *Gonta) ServeActions(w http.ResponseWriter, r *http.Request) {
	log := s.log

	var payload *slack.InteractionCallback

	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		log.Error("failed to parse payload", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	switch payload.Type {
	case slack.InteractionTypeBlockActions:
		if len(payload.ActionCallback.BlockActions) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		action := payload.ActionCallback.BlockActions[0]
		log.Debug("action.BlockID", zap.String("blockID", action.BlockID))
		w.WriteHeader(http.StatusOK)
	default:
		log.Error("Unexpected case", zap.String("payload.Type", string(payload.Type)))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Gonta) ServeCommands(w http.ResponseWriter, r *http.Request) {
	log := s.log

	cmd, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Error("Failed to parse command", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	log.Debug("Received a command", zap.String("command", cmd.Command))
	w.WriteHeader(http.StatusOK)
}

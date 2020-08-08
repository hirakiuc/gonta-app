package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/event/data"
	"github.com/hirakiuc/gonta-app/event/queue"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

var ErrUnexpectedEventType = errors.New("unexpected event type")

// Gonta describe a http server to serve gonta services.
type Gonta struct {
	log    *zap.Logger
	config *config.Config
	queue  *queue.Queue
	data   *data.Provider
}

func NewGonta(logger *zap.Logger, c *config.Config, q *queue.Queue, d *data.Provider) *Gonta {
	return &Gonta{
		log:    logger,
		config: c,
		queue:  q,
		data:   d,
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

func (s *Gonta) ServeHealth(w http.ResponseWriter, r *http.Request) {
	err := s.config.Load()
	if err != nil {
		s.log.Error("Failed to load config", zap.Error(err))
	}

	w.WriteHeader(http.StatusOK)
}

// nolint:funlen
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

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var res *slackevents.ChallengeResponse

		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Error("Failed to parse body as URLVerification event", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/plain")

		_, err = w.Write([]byte(res.Challenge))
		if err != nil {
			log.Error("Failed to write response", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

	case slackevents.CallbackEvent:
		// Enqueue the event
		s.queue.EnqueueEvent(&eventsAPIEvent)

		w.WriteHeader(http.StatusOK)

	default:
		log.Error("Unexpected event type", zap.String("type", eventsAPIEvent.Type))
		w.WriteHeader(http.StatusBadRequest)

		return
	}
}

func (s *Gonta) ServeActions(w http.ResponseWriter, r *http.Request) {
	log := s.log

	if r.Method != http.MethodPost {
		log.Debug("Invalid http method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	var payload *slack.InteractionCallback

	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		log.Error("failed to parse payload", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	log.Debug("Received a action callback", zap.String("type", string(payload.Type)))

	s.queue.EnqueueAction(payload)

	w.WriteHeader(http.StatusOK)
}

/**
 * ServeData
 *
 * This method will provide json data as external data source
 *
 * NOTE: https://api.slack.com/reference/block-kit/block-elements#external_select
 */
func (s *Gonta) ServeData(w http.ResponseWriter, r *http.Request) {
	log := s.log

	if r.Method != http.MethodPost {
		log.Debug("Invalid http request", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	// parse post parameters
	if err := r.ParseForm(); err != nil {
		log.Error("Failed to parse POST parameters", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	event, err := ParseExternalDataRequest([]byte(r.FormValue("payload")))
	if err != nil {
		log.Error("failed to parse payload", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	log.Info("Data request received", zap.String("payload", r.FormValue("payload")))

	w.Header().Set("Content-Type", "application/json")

	err = s.data.Process(event, w)
	if err != nil {
		log.Error("Failed to fetch data", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	log.Info("Finished to process a data request")
}

func (s *Gonta) ServeCommands(w http.ResponseWriter, r *http.Request) {
	log := s.log

	cmd, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Error("Failed to parse command", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	s.queue.EnqueueCommand(&cmd)

	log.Debug("Received a command", zap.String("command", cmd.Command))
	w.WriteHeader(http.StatusOK)
}

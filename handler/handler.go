package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/hirakiuc/gonta-app/log"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

// SlackEvent describe a event which sent from slack
type SlackEvent struct {
	// Challenge is a token which sent from slack on url_verification event.
	Challenge string `json:"challenge"`

	slack.AttachmentActionCallback
}

type challengeResponse struct {
	Challenge string `json:"challenge"`
}

// Serve handles the http request.
func Serve(w http.ResponseWriter, r *http.Request) {
	log := log.GetLogger()

	if r.Method != http.MethodPost {
		log.Debug("Invalid http method",
			zap.String("method", r.Method),
		)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	msg, err := parseBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if msg.Type == "url_verification" {
		respondChallenge(w, msg)
		return
	}

	// Only accept message from salck with valid token
	if msg.Token != getVerificationToken() {
		log.Debug("Invalid verification token", zap.String("verification token", msg.Token))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

	return
}

func parseBody(r *http.Request) (*SlackEvent, error) {
	log := log.GetLogger()

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read request body", zap.Error(err))
		return nil, err
	}

	jsonStr, err := url.QueryUnescape(string(buf))
	if err != nil {
		log.Error("Failed to unescape request body", zap.Error(err))
		return nil, err
	}

	var msg SlackEvent
	if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
		log.Error("Failed to decode json message from slack", zap.String("json", jsonStr))
		return nil, err
	}

	return &msg, nil
}

func getVerificationToken() string {
	return os.Getenv("VERIFICATION_TOKEN")
}

func respondChallenge(w http.ResponseWriter, msg *SlackEvent) {
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
	w.Write(res)
}

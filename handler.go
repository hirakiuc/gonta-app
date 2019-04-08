package gonta

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/hirakiuc/gonta-app/event"
	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/reply"
	"go.uber.org/zap"
)

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

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf))
	if err != nil {
		log.Error("Failed to unescape request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug("Received body", zap.String("body", jsonStr))

	parser := event.NewParser()
	// Only accept message from salck with valid token
	token, err := parser.GetToken()
	if err != nil {
		log.Error("Failed to extract token from the event", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if *token != getVerificationToken() {
		log.Debug("Invalid verification token", zap.String("verification token", token))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	eventType, err := parser.GetType(jsonStr)
	if err != nil {
		log.Error("Failed to parse the type", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch eventType {
	case "url_verification":
		e, err := parser.ParseURLVerificationEvent(jsonStr)
		if err != nil {
			log.Error("Failed to parse the URLVerificationEvent", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		replyer := reply.NewUrlVerificationReplyer()
		replyer.Reply(w)

		// Reply NewUrlVerificationReplyer
	default:
		e, err := parser.ParseCallbackEvent(jsonStr)
		if err != nil {
			log.Error("Failed to parse the CallbackEvent", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Reply
		replyer := reply.NewBeerSelectReply()
		replyer.Reply(w, e)
	}

	replyer := getReplyer(evt)
	replyer.Reply(evt)
}

func getVerificationToken() string {
	return os.Getenv("VERIFICATION_TOKEN")
}

func getReplyer(msg *event.SlackEvent) *reply.Replyer {
	switch msg.Type {
	case "app_mention":
		return reply.NewBeerSelectReply()
	case "url_verification":
		return reply.NewUrlVerificationReplyer()
	default:
		return reply.NewEmptyReplyer()
	}
}

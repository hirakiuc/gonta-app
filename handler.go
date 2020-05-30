package gonta

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/hirakiuc/gonta-app/handler"
	"github.com/hirakiuc/gonta-app/log"
	"github.com/hirakiuc/gonta-app/parser"
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

	result, err := parseBody(w, r, log)
	if err != nil {
		return
	}

	if result.Token != getVerificationToken() {
		log.Debug("Invalid verification token", zap.String("verification token", result.Token))
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	err = handleEvent(w, result, log)
	if err != nil {
		return
	}
}

func parseBody(w http.ResponseWriter, r *http.Request, log *zap.Logger) (*parser.BodyParseResult, error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return nil, err
	}

	jsonStr, err := url.QueryUnescape(string(buf))
	if err != nil {
		log.Error("Failed to unescape request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return nil, err
	}

	log.Debug("Received body", zap.String("body", jsonStr))

	bodyParser := parser.NewBodyParser()

	result, err := bodyParser.Parse(jsonStr)
	if err != nil {
		log.Error("Failed to parse event", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return nil, err
	}

	return result, nil
}

func getVerificationToken() string {
	return os.Getenv("VERIFICATION_TOKEN")
}

func handleEvent(w http.ResponseWriter, result *parser.BodyParseResult, log *zap.Logger) error {
	eventParser := parser.NewEventParser()

	switch result.Type {
	case "url_verification":
		e, err := eventParser.ParseURLVerificationEvent(result.JSON)
		if err != nil {
			log.Error("Failed to parse the URLVerificationEvent", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)

			return err
		}

		handler := handler.NewURLVerificationHandler()

		return handler.Handle(w, e)

	default:
		e, err := eventParser.ParseCallbackEvent(result.JSON)
		if err != nil {
			log.Error("Failed to parse the CallbackEvent", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)

			return err
		}

		handler := handler.NewMentionHandler()

		return handler.Handle(w, e)
	}
}

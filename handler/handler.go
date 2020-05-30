package handler

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/event"
	"go.uber.org/zap"
)

type Handler interface {
	Handle(w http.ResponseWriter, msg *event.CallbackEvent) error
}

func HandleEvent(w http.ResponseWriter, eventType string, json string, log *zap.Logger) error {
	parser := event.NewEventParser()

	switch result.Type {
	case "url_verification":
		e, err := parser.ParseURLVerificationEvent(result.JSON)
		if err != nil {
			log.Error("Failed to parse the URLVerificationEvent", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)

			return err
		}

		handler := NewURLVerificationHandler()
		return handler.Handle(w, e)

	default:
		e, err := parser.ParseCallbackEvent(result.JSON)
		if err != nil {
			log.Error("Failed to parse the CallbackEvent", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)

			return err
		}

		handler := NewMentionHandler()
		return handler.Handle(w, e)
	}

	return nil
}

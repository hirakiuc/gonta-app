package handler

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/event"
)

type Handler interface {
	Handle(w http.ResponseWriter, msg *event.CallbackEvent)
}

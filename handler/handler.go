package handler

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/model"
)

type Handler interface {
	Handle(w http.ResponseWriter, msg *model.CallbackEvent) error
}

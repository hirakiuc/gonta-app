package handler

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/model"
	"go.uber.org/zap"
)

type Handler interface {
	SetLogger(logger *zap.Logger)
	Handle(w http.ResponseWriter, msg *model.CallbackEvent) error
}

type BaseHandler struct {
	log *zap.Logger
}

func (h *BaseHandler) SetLogger(logger *zap.Logger) {
	h.log = logger
}

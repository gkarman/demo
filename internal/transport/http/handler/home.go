package handler

import (
	"log/slog"
	"net/http"

	"github.com/gkarman/demo/internal/logger"
)

type HomeHandler struct {
	log *slog.Logger
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{
	}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request, ) {
	_ = logger.FromContext(r.Context()) // если понадобится логировать
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}

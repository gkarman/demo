package handler

import (
	"log/slog"
	"net/http"
)

type HomeHandler struct {
	log *slog.Logger
}

func NewHomeHandler(log *slog.Logger) *HomeHandler {
	return &HomeHandler{
		log: log,
	}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request, ) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}

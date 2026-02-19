package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/logger"
)

type CarHandler struct {
	repo car.Repo
}

func NewCarHandler(repo car.Repo) *CarHandler {
	return &CarHandler{
		repo: repo,
	}
}

func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	cars, err := h.repo.List(r.Context())
	if err != nil {
		log.Error("get cars failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cars)
}

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/logger"
)

type CarHandler struct {
	repo car.Repository
}

func NewCarHandler(repo car.Repository) *CarHandler {
	return &CarHandler{
		repo: repo,
	}
}

func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	cars, err := h.repo.List(r.Context())
	if err != nil {
		log.Debug("get cars error", "error", err)
	}
	json.NewEncoder(w).Encode(cars)
}

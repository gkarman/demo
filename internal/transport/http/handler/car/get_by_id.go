package car

import (
	"encoding/json"
	"net/http"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/logger"
)

type GetByIdHandler struct {
	repo car.Repo
}

func NewGetByIdHandler(repo car.Repo) *GetByIdHandler {
	return &GetByIdHandler{
		repo: repo,
	}
}

func (h *GetByIdHandler) Execute(w http.ResponseWriter, r *http.Request) {
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

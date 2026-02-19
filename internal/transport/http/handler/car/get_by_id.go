package car

import (
	"encoding/json"
	"net/http"

	"github.com/gkarman/demo/internal/logger"
	car "github.com/gkarman/demo/internal/service/car"
	"github.com/go-chi/chi/v5"
)

type GetByIDHandler struct {
	service *car.GetByIDService
}

func NewGetByIDHandler(service *car.GetByIDService) *GetByIDHandler {
	return &GetByIDHandler{
		service: service,
	}
}

func (h *GetByIDHandler) Execute(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	c, err := h.service.Execute(r.Context(), id)
	if err != nil {
		log.Error("get car failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(c)
}

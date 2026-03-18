package car

import (
	"encoding/json"
	"net/http"

	"github.com/gkarman/demo/internal/logger"
	car "github.com/gkarman/demo/internal/service/car"
	"github.com/gkarman/demo/internal/service/car/requestdto"
	"github.com/go-chi/chi/v5"
)

type UpdateHandler struct {
	service *car.UpdateService
}

func NewUpdate(service *car.UpdateService) *UpdateHandler {
	return &UpdateHandler{
		service: service,
	}
}

func (h *UpdateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &requestdto.UpdateCar{
		CarId: id,
		Name:  body.Name,
	}

	resp, err := h.service.Execute(r.Context(), req)
	if err != nil {
		log.Error("update car failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

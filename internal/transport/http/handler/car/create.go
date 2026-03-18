package car

import (
	"encoding/json"
	"net/http"

	"github.com/gkarman/demo/internal/logger"
	"github.com/gkarman/demo/internal/service/car"
	"github.com/gkarman/demo/internal/service/car/requestdto"
)

type CreateHandler struct {
	service *car.CreateService
}

func NewCreate(service *car.CreateService) *CreateHandler {
	return &CreateHandler{
		service: service,
	}
}

func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")

	req := &requestdto.CreateCar{
		Name: "test",
	}

	resp, err := h.service.Execute(r.Context(), req)
	if err != nil {
		log.Error("save car failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

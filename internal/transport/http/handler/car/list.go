package car

import (
	"encoding/json"
	"net/http"

	"github.com/gkarman/demo/internal/logger"
	"github.com/gkarman/demo/internal/service/car"
	"github.com/gkarman/demo/internal/service/car/requestdto"
)

type ListHandler struct {
	service *car.List
}

func NewList(service *car.List) *ListHandler {
	return &ListHandler{
		service: service,
	}
}

func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	resp, err := h.service.Execute(r.Context(), requestdto.GetList{})
	if err != nil {
		log.Error("get cars failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

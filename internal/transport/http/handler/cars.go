package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarHandler struct {
	db  *pgxpool.Pool
}

func NewCarHandler(db *pgxpool.Pool, ) *CarHandler {
	return &CarHandler{
		db:  db,
	}
}

func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())
	_ = log
	cars := car.Car{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarHandler struct {
	log *slog.Logger
	db  *pgxpool.Pool
}

func NewCarHandler(log *slog.Logger, db *pgxpool.Pool, ) *CarHandler {
	return &CarHandler{
		log: log,
		db:  db,
	}
}

func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	cars := car.Car{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

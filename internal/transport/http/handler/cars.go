package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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

type Car struct {
	ID   string    `json:"id"`
	Name string `json:"name"`
}

func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second, )
	defer cancel()

	rows, err := h.db.Query(
		ctx,
		"SELECT id, name FROM cars",
	)

	if err != nil {
		h.log.Error("query failed", "error", err)
		http.Error(
			w,
			"internal error",
			http.StatusInternalServerError,
		)
		return
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var c Car
		if err := rows.Scan(&c.ID, &c.Name, ); err != nil {
			h.log.Error("scan failed", "error", err)
			continue
		}
		cars = append(cars, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

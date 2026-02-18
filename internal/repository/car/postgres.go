package car

import (
	"context"
	"errors"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(ctx context.Context) ([]*car.Car, error){
	return []*car.Car{}, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*car.Car, error) {
	const q = `
		SELECT id, name
		FROM cars
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, q, id)

	var c car.Car
	if err := row.Scan(&c.ID, &c.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}


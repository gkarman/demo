package car

import (
	"context"
	"errors"

	"github.com/gkarman/demo/internal/domain/car"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) GetByID(ctx context.Context, id string) (*car.Car, error) {
	const q = `
		SELECT id, name
		FROM users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, q, id)

	var u car.Car
	if err := row.Scan(&u.ID, &u.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

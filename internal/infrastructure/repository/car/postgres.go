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

func (r *Repository) List(ctx context.Context) ([]*car.Car, error) {
	const q = `
		SELECT id, name
		FROM cars
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	carList := make([]*car.Car, 0)
	for rows.Next() {
		var c car.Car
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		carList = append(carList, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return carList, nil
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
			return nil, car.ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}

func (r *Repository) Save(ctx context.Context, car *car.Car) error {
	const q = `
		INSERT INTO cars (id, name) VALUES ($1, $2)
	`

	_, err := r.db.Exec(
		ctx,
		q,
		car.ID,
		car.Name,
	)

	return err
}

func (r *Repository) Update(ctx context.Context, car *car.Car) error {
	const q = `
		UPDATE cars SET name = $2 WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		q,
		car.ID,
		car.Name,
	)

	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	const q = `
		DELETE FROM cars WHERE id = $1
	`

	result, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return car.ErrNotFound
	}

	return nil
}

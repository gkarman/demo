package car

import "context"

type Repository interface {
	List(ctx context.Context) ([]*Car, error)
	GetByID(ctx context.Context, id string) (*Car, error)
}
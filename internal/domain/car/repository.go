package car

import "context"

type Repo interface {
	List(ctx context.Context) ([]*Car, error)
	GetByID(ctx context.Context, id string) (*Car, error)
}
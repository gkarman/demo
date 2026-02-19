package car

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/domain/car"
)

type GetByIDService struct {
	repo car.Repo
}

func NewGetByIDService(repo car.Repo) *GetByIDService {
	return &GetByIDService{
		repo: repo,
	}
}

func (s *GetByIDService) Execute(ctx context.Context, id string) (*car.Car, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(`GetByIDService.Execute: %w`, err)
	}
	return c, nil
}
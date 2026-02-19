package car

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/domain/car"
)

type GetByIdService struct {
	repo car.Repo
}

func NewGetByIdService(repo car.Repo) *ListService {
	return &ListService{
		repo: repo,
	}
}

func (s *GetByIdService) Execute(ctx context.Context, id string) (*car.Car, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(`GetByIdService.Execute: %w`, err)
	}
	return c, nil
}
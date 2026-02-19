package car

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/domain/car"
)

type ListService struct {
	repo car.Repo
}

func NewListService(repo car.Repo) *ListService {
	return &ListService{
		repo: repo,
	}
}

func (s *ListService) Execute(ctx context.Context) ([]*car.Car, error) {
	c, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf(`ListService.Execute: %w`, err)
	}
	return c, nil
}
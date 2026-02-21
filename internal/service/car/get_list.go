package car

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/service/car/mapper"
	"github.com/gkarman/demo/internal/service/car/requestdto"
	"github.com/gkarman/demo/internal/service/car/responsedto"
)

type ListService struct {
	repo car.Repo
}

func NewListService(repo car.Repo) *ListService {
	return &ListService{
		repo: repo,
	}
}

func (s *ListService) Execute(ctx context.Context, _ requestdto.GetList) (*responsedto.GetList, error) {
	cars, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf(`ListService.handel: %w`, err)
	}
	resp := &responsedto.GetList{
		Cars: mapper.CarsFromDomain(cars),
	}
	return resp, nil
}

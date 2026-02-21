package car

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/service/car/mapper"
	"github.com/gkarman/demo/internal/service/car/requestdto"
	"github.com/gkarman/demo/internal/service/car/responsedto"
)

type GetCarService struct {
	repo car.Repo
}

func NewGetCarService(repo car.Repo) *GetCarService {
	return &GetCarService{
		repo: repo,
	}
}

func (s *GetCarService) Execute(ctx context.Context, req *requestdto.GetCar) (*responsedto.GetCar, error) {
	c, err := s.repo.GetByID(ctx, req.CarId)
	if err != nil {
		return nil, fmt.Errorf(`GetCarService.handel: %w`, err)
	}
	resp := &responsedto.GetCar{
		Car: mapper.CarFromDomain(c),
	}
	return resp, nil
}

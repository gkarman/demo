package car

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/service/car/requestdto"
	"github.com/gkarman/demo/internal/service/car/responsedto"
	"github.com/google/uuid"
)

type CreateService struct {
	repo car.Repo
}

func NewCreate(repo car.Repo) *CreateService {
	return &CreateService{
		repo: repo,
	}
}

func (s *CreateService) Execute(ctx context.Context, req *requestdto.CreateCar) (*responsedto.CreateCar, error) {
	if req.Name == "" {
		return nil, car.ErrEmptyName
	}

	id := uuid.New()
	c := car.New(id.String(), req.Name)

	if err := s.repo.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("CreateService.Execute: %w", err)
	}

	return &responsedto.CreateCar{
		ID: c.ID,
	}, nil
}

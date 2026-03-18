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
	id := uuid.New()
	name := req.Name

	c := car.New(id.String(), name)
	err := s.repo.Save(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("failed to save car in CreateService: %w", err)
	}

	resp := &responsedto.CreateCar{
		ID: c.ID,
	}

	return resp, nil
}

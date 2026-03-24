package car

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/infrastructure/eventbus"
	"github.com/gkarman/demo/internal/service/car/requestdto"
	"github.com/gkarman/demo/internal/service/car/responsedto"
	"github.com/google/uuid"
)

type CreateService struct {
	repo     car.Repo
	eventBus *eventbus.EventBus
}

func NewCreate(repo car.Repo, bus *eventbus.EventBus) *CreateService {
	return &CreateService{
		repo:     repo,
		eventBus: bus,
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

	for _, event := range c.PullEvents() {
		s.eventBus.Publish(ctx, event)
	}

	return &responsedto.CreateCar{
		ID: c.ID,
	}, nil
}

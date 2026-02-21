package mapper

import (
	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/service/car/responsedto"
)

func CarFromDomain(c *car.Car) *responsedto.Car {
	if c == nil {
		return nil
	}
	return &responsedto.Car{
		ID:   c.ID,
		Name: c.Name,
	}
}

func CarsFromDomain(cars []*car.Car) []*responsedto.Car {
	if cars == nil {
		return nil
	}
	out := make([]*responsedto.Car, len(cars))
	for i, c := range cars {
		out[i] = CarFromDomain(c)
	}
	return out
}

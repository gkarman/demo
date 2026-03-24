package car

import (
	"github.com/gkarman/demo/internal/domain/car/events"
)

type Car struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	events []events.Domain
}

func New(id string, name string) *Car {
	car := &Car{
		ID:   id,
		Name: name,
	}

	e := events.NewCarCreated(car.ID, car.Name)
	car.addEvent(e)
	return car
}

func (c *Car) addEvent(e events.Domain) {
	c.events = append(c.events, e)
}

func (c *Car) PullEvents() []events.Domain {
	evs := c.events
	c.events = nil
	return evs
}

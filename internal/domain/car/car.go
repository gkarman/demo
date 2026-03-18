package car

type Car struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func New(id string, name string) *Car {
	car := &Car{
		ID:   id,
		Name: name,
	}
	return car
}

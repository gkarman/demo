package events

import (
	"log/slog"
	"reflect"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/car/handlers"
	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
)

func RegisterEventHandlers(d *dispatcher.Dispatcher, log *slog.Logger, publisher application.Publisher) {
	d.Register(reflect.TypeOf(&car.Created{}), handlers.CarCreatedLogHandler(log))
	d.Register(reflect.TypeOf(&car.Created{}), handlers.CarCreatedToRabbitHandler(publisher, log))
}

package events

import (
	"log/slog"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/car/handlers"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
)

func RegisterEventHandlers(d *dispatcher.Dispatcher, log *slog.Logger, publisher application.Publisher) {
	d.Register("car.created", handlers.CarCreatedLogHandler(log))
	d.Register("car.created", handlers.CarCreatedToRabbitHandler(publisher, log))
}

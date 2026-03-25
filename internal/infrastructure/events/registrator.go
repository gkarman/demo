package events

import (
	"log/slog"

	"github.com/gkarman/demo/internal/application/service/car/handlers"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/mq"
)

func RegisterEventHandlers(d *dispatcher.Dispatcher, log *slog.Logger, publisher mq.Publisher) {
	d.Register("car.created", handlers.CarCreatedLogHandler(log))
	d.Register("car.created", handlers.CarCreatedToRabbitHandler(publisher, log))
}

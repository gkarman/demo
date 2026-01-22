package main

import (
	"log/slog"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/logger"
)

func main() {
	cfg := config.MustLoad("configs/config.yaml")
	log := logger.New(logger.Config{Level: cfg.Logger.Level})
	slog.SetDefault(log)

}

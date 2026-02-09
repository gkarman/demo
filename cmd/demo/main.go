package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gkarman/demo/internal/db"
	"github.com/gkarman/demo/internal/repository/user"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/logger"
)

type User struct {
	ID   string
	Name string
}

func main() {
	cfg := config.MustLoad("configs/config.yaml")
	log := logger.New(logger.Config{Level: cfg.Logger.Level})
	slog.SetDefault(log)
	log.Info("config", cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pg, err := db.New(ctx, db.Config{
		DSN:             cfg.DB.DSN(),
		MaxConns:        cfg.DB.MaxConnections,
		MinConns:        cfg.DB.MinConnections,
		MaxConnLifetime: time.Duration(cfg.DB.MaxConnectionLifeTimeMinutes) * time.Minute,
		MaxConnIdleTime: time.Duration(cfg.DB.MaxConnectionIdleTimeMinutes) * time.Minute,
	})

	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pg.Close()

	log.Info("database connected successfully")



	userRepo := user.New(pg.Pool)
	u, err := userRepo.GetByID(ctx, "d1c38fc3-6b21-4b62-bf4c-5a86c05235a2")
	if err != nil {
		log.Error("user repo", err)
	}
	log.Info("user", u)
}

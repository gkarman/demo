package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gkarman/demo/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	command := os.Args[1]
	fs := flag.NewFlagSet(command, flag.ExitOnError)
	step := fs.Int("step", 0, "number of migrations to apply/rollback")

	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	if *step < 0 {
		log.Fatal("--step must be a positive number")
	}

	cfg := config.MustLoad("configs/config.yaml")
	dsn := buildPostgresDSN(cfg.DB)

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("failed to init migrate: %v", err)
	}

	switch command {
	case "up":
		runUp(m, *step)
	case "down":
		runDown(m, *step)
	case "version":
		runVersion(m)
	default:
		usage()
		os.Exit(1)
	}
}

func buildPostgresDSN(db config.DBConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pass,
		db.Host,
		db.Port,
		db.Name,
	)
}

func runUp(m *migrate.Migrate, step int) {
	var err error

	if step > 0 {
		err = m.Steps(step)
	} else {
		err = m.Up()
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration up failed: %v", err)
	}

	log.Println("migrations applied")
}

func runDown(m *migrate.Migrate, step int) {
	var err error

	if step > 0 {
		err = m.Steps(-step)
	} else {
		err = m.Down()
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration down failed: %v", err)
	}

	log.Println("migrations rolled back")
}

func runVersion(m *migrate.Migrate) {
	version, dirty, err := m.Version()
	if err != nil {
		log.Fatalf("failed to get version: %v", err)
	}

	fmt.Printf("version=%d dirty=%v\n", version, dirty)
}

func usage() {
	fmt.Println(`Usage:
  migrate up [--step=N]       apply all or N migrations
  migrate down [--step=N]     rollback all or N migrations
  migrate version             print current version`)
}

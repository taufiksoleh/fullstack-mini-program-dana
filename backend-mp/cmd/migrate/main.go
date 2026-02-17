package main

import (
	"log/slog"
	"os"

	pgmodel "github.com/taufiksoleh/backend-mp/internal/repository/postgres/model"
	"github.com/taufiksoleh/backend-mp/pkg/config"
	"github.com/taufiksoleh/backend-mp/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		slog.Error("migrate: failed to connect to database", "err", err)
		os.Exit(1)
	}

	slog.Info("migrate: running AutoMigrate...")

	if err := db.AutoMigrate(
		&pgmodel.Facility{},
		&pgmodel.Property{},
		&pgmodel.PropertyImage{},
	); err != nil {
		slog.Error("migrate: failed to run migrations", "err", err)
		os.Exit(1)
	}

	slog.Info("migrate: migrations completed successfully")
}

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	deliveryhttp "github.com/taufiksoleh/backend-mp/internal/delivery/http"
	"github.com/taufiksoleh/backend-mp/internal/delivery/http/handler"
	postgresrepo "github.com/taufiksoleh/backend-mp/internal/repository/postgres"
	"github.com/taufiksoleh/backend-mp/internal/usecase"
	"github.com/taufiksoleh/backend-mp/pkg/config"
	"github.com/taufiksoleh/backend-mp/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	propertyRepo := postgresrepo.NewPropertyRepository(db)
	propertyUsecase := usecase.NewPropertyUsecase(propertyRepo, db)
	propertyHandler := handler.NewPropertyHandler(propertyUsecase)

	r := deliveryhttp.NewRouter(propertyHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		slog.Info("server starting", "port", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

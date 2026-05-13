package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/config"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/server"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
)

func main() {
	production := os.Getenv("APP_ENV") == "production"

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if production {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)

	if err := run(production); err != nil {
		slog.Error("fatal error", "error", err)

		os.Exit(1)
	}
}

func run(production bool) error {
	cfg, err := config.Load()

	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx := context.Background()

	db, err := store.Open(ctx, cfg.DatabaseURL)

	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	if err := db.EnsureAdmin(ctx, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		return fmt.Errorf("ensure admin: %w", err)
	}

	app, err := server.New(db, production)

	if err != nil {
		return fmt.Errorf("create server: %w", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.Routes(),
	}

	slog.Info("go-fetch analytics listening", "addr", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("listen http: %w", err)
	}

	return nil
}

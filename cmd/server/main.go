package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/config"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/server"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cfg config.Config) error {
	if err := store.Migrate(ctx, cfg.DatabaseURL); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	db, err := store.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	if err := db.EnsureAdmin(ctx, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		return fmt.Errorf("ensure admin: %w", err)
	}

	app := server.New(db)

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: app.Routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("listen http: %w", err)
	}

	return nil
}

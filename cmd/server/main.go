package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/config"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/database"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/router"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := run(ctx, appConfig); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, appConfig config.Config) error {
	if err := database.Migrate(ctx, appConfig.DatabaseURL); err != nil {
		return fmt.Errorf("run database migrations: %w", err)
	}

	databasePool, err := database.Open(ctx, appConfig.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open database connection: %w", err)
	}
	defer databasePool.Close()

	dataStore := repository.New(databasePool)

	if err := dataStore.EnsureAdminUser(ctx, appConfig.AdminUsername, appConfig.AdminPassword); err != nil {
		return fmt.Errorf("initialize admin user: %w", err)
	}

	application := router.New(dataStore)
	httpServer := &http.Server{
		Addr:         appConfig.ListenAddr,
		Handler:      application.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("start HTTP server: %w", err)
	}

	return nil
}

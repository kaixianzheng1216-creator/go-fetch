package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/config"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/database"
)

func main() {
	databaseURL, err := config.LoadDatabaseURL()
	if err != nil {
		slog.Error("load database config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := database.Migrate(ctx, databaseURL); err != nil {
		slog.Error("run database migrations", "error", err)
		os.Exit(1)
	}
}

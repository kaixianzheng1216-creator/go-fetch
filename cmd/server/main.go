package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/app"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/config"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, appConfig); err != nil {
		slog.Error("run application", "error", err)
		os.Exit(1)
	}
}

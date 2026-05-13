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

	"github.com/kaixianzheng1216-creator/go-fetch/internal/config"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/server"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
)

func main() {
	production := os.Getenv("APP_ENV") == "production"

	logger := newLogger(production)
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		fatal(logger, err)
	}

	ctx := context.Background()

	db, err := store.Open(ctx, cfg.DatabaseURL)

	if err != nil {
		fatal(logger, err)
	}

	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		fatal(logger, err)
	}

	if err := db.EnsureAdmin(ctx, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		fatal(logger, err)
	}

	app, err := server.New(db, production)
	if err != nil {
		fatal(logger, err)
	}

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           app.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
	}

	go func() {
		logger.Info(
			"go-fetch analytics listening",
			"addr",
			httpServer.Addr,
		)

		if err := httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			fatal(logger, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
	}
}

func newLogger(production bool) *slog.Logger {
	if production {
		return slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		)
	}

	return slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	)
}

func fatal(logger *slog.Logger, err error) {
	logger.Error("fatal error", "error", err)
	os.Exit(1)
}

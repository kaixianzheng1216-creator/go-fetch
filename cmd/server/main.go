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

	"go-fetch/internal/config"
	"go-fetch/internal/server"
	"go-fetch/internal/store"
)

func main() {
	logger := newLogger(os.Getenv("APP_ENV"))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		fatal(logger, err)
	}
	logger = newLogger(cfg.Environment)
	slog.SetDefault(logger)

	ctx := context.Background()
	db, err := store.Open(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
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

	app, err := server.New(cfg, db)
	if err != nil {
		fatal(logger, err)
	}

	httpServer := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           app.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	go func() {
		logger.Info("go-fetch analytics listening", "addr", cfg.ListenAddr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fatal(logger, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
	}
}

func newLogger(environment string) *slog.Logger {
	if environment == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func fatal(logger *slog.Logger, err error) {
	logger.Error("fatal error", "error", err)
	os.Exit(1)
}

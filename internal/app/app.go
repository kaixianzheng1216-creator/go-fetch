// Package app wires dependencies and runs the HTTP server.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/config"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/database"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

func Run(ctx context.Context, appConfig config.Config) error {
	if err := database.Migrate(ctx, appConfig.DatabaseURL); err != nil {
		return fmt.Errorf("run database migrations: %w", err)
	}

	dbPool, err := database.Open(ctx, appConfig.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open database connection: %w", err)
	}
	defer dbPool.Close()

	dataStore := repository.New(dbPool)
	users := service.NewUserService(dataStore)
	if err := users.EnsureAdminUser(ctx, appConfig.AdminUsername, appConfig.AdminPassword); err != nil {
		return fmt.Errorf("initialize admin user: %w", err)
	}

	sessionManager := session.NewManager(dbPool, session.Config{
		CookieSecure: appConfig.SessionCookieSecure,
		Lifetime:     appConfig.SessionLifetime,
	})

	httpServer := &http.Server{
		Addr: appConfig.ListenAddr,
		Handler: httpapi.New(httpapi.Services{
			Auth:       service.NewAuthService(dataStore),
			Collection: service.NewCollectionService(dataStore),
			Stats:      service.NewStatsService(dataStore),
			Users:      users,
			Websites:   service.NewWebsiteService(dataStore),
		}, sessionManager, httpapi.Config{
			CollectCORSAllowedOrigins: appConfig.CollectCORSAllowedOrigins,
			RequestTimeout:            appConfig.HTTPRequestTimeout,
			TrustProxyHeaders:         appConfig.TrustProxyHeaders,
		}),
		ReadTimeout:  appConfig.HTTPReadTimeout,
		WriteTimeout: appConfig.HTTPWriteTimeout,
		IdleTimeout:  appConfig.HTTPIdleTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("starting HTTP server", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
			return
		}
		serverErrors <- nil
	}()

	select {
	case err := <-serverErrors:
		if err != nil {
			return fmt.Errorf("start HTTP server: %w", err)
		}
		return nil
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), appConfig.HTTPShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		select {
		case err := <-serverErrors:
			if err != nil {
				return fmt.Errorf("stop HTTP server: %w", err)
			}
		case <-time.After(time.Second):
		}

		return nil
	}
}

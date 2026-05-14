package store

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, databaseURL string) error {
	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("create migration database handle: %w", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Warn("close migration database", "error", err)
		}
	}()

	sqlDB.SetMaxOpenConns(1)

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping migration database: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, migrations.FS)
	if err != nil {
		return fmt.Errorf("create migration provider: %w", err)
	}

	results, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	for _, result := range results {
		slog.Info(
			"migration applied",
			"version", result.Source.Version,
			"path", result.Source.Path,
			"direction", result.Direction,
			"duration", result.Duration,
		)
	}

	slog.Info("migrations complete", "applied", len(results))

	return nil
}

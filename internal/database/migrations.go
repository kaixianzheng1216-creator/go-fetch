package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kaixianzheng1216-creator/go-fetch/migrations"

	_ "github.com/jackc/pgx/v5/stdlib" // Registers the pgx database/sql driver for goose.
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, databaseURL string) (err error) {
	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open migration database handle: %w", err)
	}
	defer func() {
		if closeErr := sqlDB.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close migration database handle: %w", closeErr))
		}
	}()

	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, migrations.FS)
	if err != nil {
		return fmt.Errorf("create migration provider: %w", err)
	}

	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("run database migrations: %w", err)
	}

	return nil
}

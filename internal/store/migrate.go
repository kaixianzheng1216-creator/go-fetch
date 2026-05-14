package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, databaseURL string) error {
	sqlDB, err := sql.Open("pgx", databaseURL)

	if err != nil {
		return fmt.Errorf("open migration database: %w", err)
	}

	defer sqlDB.Close()

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}

	if err := goose.UpContext(ctx, sqlDB, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

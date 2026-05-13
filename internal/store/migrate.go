package store

import (
	"context"
	"database/sql"

	"go-fetch/internal/store/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (s *Store) Migrate(ctx context.Context) error {
	sqlDB, err := sql.Open("pgx", s.databaseURL)
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	sqlDB.SetMaxOpenConns(int(s.maxConns))
	sqlDB.SetMaxIdleConns(int(s.maxConns))

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.UpContext(ctx, sqlDB, ".")
}

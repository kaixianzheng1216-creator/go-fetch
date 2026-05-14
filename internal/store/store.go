package store

import (
	"context"
	"fmt"

	storedb "github.com/kaixianzheng1216-creator/go-fetch/internal/store/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db      *pgxpool.Pool
	queries *storedb.Queries
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Store{db: pool, queries: storedb.New(pool)}, nil
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *Store) Pool() *pgxpool.Pool {
	return s.db
}

func (s *Store) Ping(ctx context.Context) error {
	if err := s.db.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

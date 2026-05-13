package store

import (
	"context"

	storedb "github.com/kaixianzheng1216-creator/go-fetch/internal/store/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db          *pgxpool.Pool
	queries     *storedb.Queries
	databaseURL string
	maxConns    int32
}

func Open(ctx context.Context, databaseURL string, maxConns int32) (*Store, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	if maxConns > 0 {
		poolConfig.MaxConns = maxConns
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Store{db: pool, queries: storedb.New(pool), databaseURL: databaseURL, maxConns: poolConfig.MaxConns}, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Pool() *pgxpool.Pool {
	return s.db
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

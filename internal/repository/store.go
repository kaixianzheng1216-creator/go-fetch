package repository

import (
	"context"
	"fmt"

	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	databasePool *pgxpool.Pool
	queries      *storesqlc.Queries
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("解析数据库 URL 失败: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("创建数据库连接池失败: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return &Store{databasePool: pool, queries: storesqlc.New(pool)}, nil
}

func (store *Store) Close() {
	if store.databasePool != nil {
		store.databasePool.Close()
	}
}

func (store *Store) Pool() *pgxpool.Pool {
	return store.databasePool
}

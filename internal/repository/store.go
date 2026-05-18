package repository

import (
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool    *pgxpool.Pool
	queries *storesqlc.Queries
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool, queries: storesqlc.New(pool)}
}

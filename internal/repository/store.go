package repository

import (
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	databasePool *pgxpool.Pool
	queries      *storesqlc.Queries
}

func New(databasePool *pgxpool.Pool) *Store {
	return &Store{databasePool: databasePool, queries: storesqlc.New(databasePool)}
}

func (store *Store) Close() {
	if store.databasePool != nil {
		store.databasePool.Close()
	}
}

func (store *Store) Pool() *pgxpool.Pool {
	return store.databasePool
}

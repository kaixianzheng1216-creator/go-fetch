package repository

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 generate -f ../../sqlc.yaml

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"
)

type Store struct {
	pool    *pgxpool.Pool
	queries *storesqlc.Queries
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool, queries: storesqlc.New(pool)}
}

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func toUser(id uuid.UUID, username, passwordHash string, createdAt time.Time, updatedAt, deletedAt *time.Time) domain.User {
	return domain.User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		DeletedAt:    deletedAt,
	}
}

func toWebsite(id uuid.UUID, name, domainName string, createdAt time.Time) domain.Website {
	return domain.Website{
		ID:        id,
		Name:      name,
		Domain:    domainName,
		CreatedAt: createdAt,
	}
}

package repository

import (
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toUser(id uuid.UUID, username, passwordHash string, createdAt time.Time, updatedAt, deletedAt pgtype.Timestamptz) domain.User {
	return domain.User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    optionalTimeFrom(updatedAt),
		DeletedAt:    optionalTimeFrom(deletedAt),
	}
}

func toWebsite(id uuid.UUID, name, domainName string, createdAt time.Time) domain.Website {
	return domain.Website{
		ID:         id,
		Name:       name,
		DomainName: domainName,
		CreatedAt:  createdAt,
	}
}

func pgFloat8(value *float64) pgtype.Float8 {
	if value == nil {
		return pgtype.Float8{}
	}

	return pgtype.Float8{
		Float64: *value,
		Valid:   true,
	}
}

func pgOptionalTime(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}

	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func optionalTimeFrom(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

package store

import (
	"math"
	"math/big"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toUser(id uuid.UUID, username, passwordHash string, createdAt time.Time, updatedAt, deletedAt pgtype.Timestamptz) domain.User {
	return domain.User{
		ID:           id.String(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    optionalTimeFrom(updatedAt),
		DeletedAt:    optionalTimeFrom(deletedAt),
	}
}

func toWebsite(id uuid.UUID, name, domainName string, createdAt time.Time) domain.Website {
	return domain.Website{
		ID:        id.String(),
		Name:      name,
		Domain:    domainName,
		CreatedAt: createdAt,
	}
}

func pgNumeric(value *float64) pgtype.Numeric {
	if value == nil {
		return pgtype.Numeric{}
	}

	scaled := math.Round(*value * 10000)

	return pgtype.Numeric{
		Int:   big.NewInt(int64(scaled)),
		Exp:   -4,
		Valid: true,
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

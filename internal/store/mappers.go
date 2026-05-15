package store

import (
	"math"
	"math/big"
	"time"

	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/user"
	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/website"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toUser(id uuid.UUID, username, passwordHash string, createdAt time.Time, updatedAt, deletedAt pgtype.Timestamptz) userdomain.User {
	return userdomain.User{
		ID:           id.String(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    optionalTimeFrom(updatedAt),
		DeletedAt:    optionalTimeFrom(deletedAt),
	}
}

func toWebsite(id uuid.UUID, name, domainName string, createdAt time.Time) websitedomain.Website {
	return websitedomain.Website{
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

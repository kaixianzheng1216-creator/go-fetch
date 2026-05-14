package store

import (
	"math"
	"math/big"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toUser(userUUID uuid.UUID, username, passwordHash, logoURL, displayName string, createdAt, updatedAt, deletedAt pgtype.Timestamptz) domain.User {
	return domain.User{
		ID:           userUUID.String(),
		Username:     username,
		PasswordHash: passwordHash,
		LogoURL:      logoURL,
		DisplayName:  displayName,
		CreatedAt:    timeFrom(createdAt),
		UpdatedAt:    optionalTimeFrom(updatedAt),
		DeletedAt:    optionalTimeFrom(deletedAt),
	}
}

func toWebsite(websiteUUID uuid.UUID, name, websiteDomain string, createdAt pgtype.Timestamptz) domain.Website {
	return domain.Website{
		ID:        websiteUUID.String(),
		Name:      name,
		Domain:    websiteDomain,
		CreatedAt: timeFrom(createdAt),
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

func pgTime(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: value, Valid: true}
}

func pgOptionalTime(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}

	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func timeFrom(value pgtype.Timestamptz) time.Time {
	if !value.Valid {
		return time.Time{}
	}

	return value.Time
}

func optionalTimeFrom(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

package store

import (
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toUser(userUUID uuid.UUID, username, passwordHash string, createdAt time.Time) domain.User {
	return domain.User{
		ID:           userUUID.String(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
	}
}

func toWebsite(websiteUUID uuid.UUID, name, websiteDomain string, createdAt time.Time) domain.Website {
	return domain.Website{
		ID:        websiteUUID.String(),
		Name:      name,
		Domain:    websiteDomain,
		CreatedAt: createdAt,
	}
}

func pgFloat(value *float64) pgtype.Float8 {
	if value == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *value, Valid: true}
}

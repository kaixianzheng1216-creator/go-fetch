package user

import (
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain/shared"
)

type User struct {
	ID           shared.ID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
}

type Website struct {
	ID        uuid.UUID
	Name      string
	Domain    string
	CreatedAt time.Time
}

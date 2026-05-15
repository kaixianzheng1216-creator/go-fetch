package domain

import (
	"time"

	"github.com/google/uuid"
)

type Website struct {
	ID        uuid.UUID
	Name      string
	Domain    string
	CreatedAt time.Time
}

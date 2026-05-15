package website

import (
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain/shared"
)

type Website struct {
	ID        shared.ID
	Name      string
	Domain    string
	CreatedAt time.Time
}

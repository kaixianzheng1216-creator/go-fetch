package httpapi

import (
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

type User struct {
	ID        string     `json:"id" format:"uuid"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" required:"true" minLength:"1"`
	Password string `json:"password" required:"true" minLength:"1" writeOnly:"true"`
}

type LoginResponse struct {
	User User `json:"user"`
}

func UserFromDomain(user domain.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

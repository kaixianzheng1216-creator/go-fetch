package response

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

type LoginResponse struct {
	User User `json:"user"`
}

type LoginOutput struct {
	Body LoginResponse
}

type UserOutput struct {
	Body User
}

func NewLoginOutput(response LoginResponse) *LoginOutput {
	return &LoginOutput{Body: response}
}

func NewUserOutput(user User) *UserOutput {
	return &UserOutput{Body: user}
}

func ToUser(user domain.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

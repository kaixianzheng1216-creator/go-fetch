package user

import (
	"time"

	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/user"
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

type loginOutput struct {
	Body LoginResponse
}

type userOutput struct {
	Body User
}

type UserOK struct {
	OK bool `json:"ok"`
}

type okOutput struct {
	Body UserOK
}

func newLoginOutput(response LoginResponse) *loginOutput {
	return &loginOutput{Body: response}
}

func newUserOutput(user User) *userOutput {
	return &userOutput{Body: user}
}

func newOKOutput() *okOutput {
	return &okOutput{Body: UserOK{OK: true}}
}

func ToUser(user userdomain.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

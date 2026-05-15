package auth

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"

	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/user"
)

type Store interface {
	GetUserByUsername(ctx context.Context, username string) (userdomain.User, error)
}

type Sessions interface {
	RenewToken(ctx context.Context) error
	Put(ctx context.Context, key string, val any)
	Destroy(ctx context.Context) error
}

type Handler struct {
	store       Store
	sessions    Sessions
	userIDKey   string
	currentUser func(context.Context) userdomain.User
	isNotFound  func(error) bool
}

func New(
	dataStore Store,
	sessions Sessions,
	userIDKey string,
	currentUser func(context.Context) userdomain.User,
	isNotFound func(error) bool,
) Handler {
	return Handler{
		store:       dataStore,
		sessions:    sessions,
		userIDKey:   userIDKey,
		currentUser: currentUser,
		isNotFound:  isNotFound,
	}
}

type loginRequest struct {
	Body LoginRequest
}

type emptyRequest struct{}

func (h Handler) Login(ctx context.Context, request *loginRequest) (*loginOutput, error) {
	user, err := h.store.GetUserByUsername(ctx, request.Body.Username)
	if err != nil {
		if h.isNotFound(err) {
			return nil, huma.Error401Unauthorized("用户名或密码错误")
		}

		return nil, huma.Error500InternalServerError("加载用户失败")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Body.Password)) != nil {
		return nil, huma.Error401Unauthorized("用户名或密码错误")
	}

	if err := h.startSession(ctx, user.ID); err != nil {
		return nil, huma.Error500InternalServerError("创建登录会话失败")
	}

	response := LoginResponse{
		User: ToUser(user),
	}

	return newLoginOutput(response), nil
}

func (h Handler) Logout(ctx context.Context, _ *emptyRequest) (*okOutput, error) {
	if err := h.sessions.Destroy(ctx); err != nil {
		return nil, huma.Error500InternalServerError("退出登录失败")
	}

	return newOKOutput(), nil
}

func (h Handler) Me(ctx context.Context, _ *emptyRequest) (*userOutput, error) {
	response := ToUser(h.currentUser(ctx))

	return newUserOutput(response), nil
}

func (h Handler) startSession(ctx context.Context, userID string) error {
	if err := h.sessions.RenewToken(ctx); err != nil {
		return fmt.Errorf("刷新会话令牌失败: %w", err)
	}

	h.sessions.Put(ctx, h.userIDKey, userID)

	return nil
}

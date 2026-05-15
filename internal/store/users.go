package store

import (
	"context"
	"fmt"

	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/store/sqlc"
	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (store *Store) EnsureAdminUser(ctx context.Context, username, password string) error {
	count, err := store.queries.CountUsers(ctx)
	if err != nil {
		return fmt.Errorf("统计用户数量失败: %w", err)
	}

	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成管理员密码哈希失败: %w", err)
	}

	if err := store.queries.CreateUser(ctx, storesqlc.CreateUserParams{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hash),
	}); err != nil {
		return fmt.Errorf("创建管理员用户失败: %w", err)
	}

	return nil
}

func (store *Store) GetUserByUsername(ctx context.Context, username string) (userdomain.User, error) {
	row, err := store.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return userdomain.User{}, fmt.Errorf("按用户名查询用户失败: %w", mapNotFound(err))
	}

	return toUser(row.ID, row.Username, row.PasswordHash, row.CreatedAt, row.UpdatedAt, row.DeletedAt), nil
}

func (store *Store) GetUserByID(ctx context.Context, userID string) (userdomain.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return userdomain.User{}, fmt.Errorf("解析用户 ID 失败: %w", err)
	}

	row, err := store.queries.GetUserByID(ctx, userUUID)
	if err != nil {
		return userdomain.User{}, fmt.Errorf("按 ID 查询用户失败: %w", mapNotFound(err))
	}

	return toUser(row.ID, row.Username, row.PasswordHash, row.CreatedAt, row.UpdatedAt, row.DeletedAt), nil
}

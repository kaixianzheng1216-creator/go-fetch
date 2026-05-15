package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, databaseURL string) error {
	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("打开迁移数据库句柄失败: %w", err)
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, migrations.FS)
	if err != nil {
		return fmt.Errorf("创建迁移执行器失败: %w", err)
	}

	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("执行数据库迁移失败: %w", err)
	}

	return nil
}

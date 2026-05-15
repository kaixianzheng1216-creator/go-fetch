package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/config"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/server"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := run(ctx, appConfig); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, appConfig config.Config) error {
	if err := store.Migrate(ctx, appConfig.DatabaseURL); err != nil {
		return fmt.Errorf("执行数据库迁移失败: %w", err)
	}

	dataStore, err := store.Open(ctx, appConfig.DatabaseURL)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %w", err)
	}
	defer dataStore.Close()

	if err := dataStore.EnsureAdminUser(ctx, appConfig.AdminUsername, appConfig.AdminPassword); err != nil {
		return fmt.Errorf("初始化管理员用户失败: %w", err)
	}

	app := server.New(dataStore)

	httpServer := &http.Server{
		Addr:         appConfig.ListenAddr,
		Handler:      app.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("启动 HTTP 服务失败: %w", err)
	}

	return nil
}

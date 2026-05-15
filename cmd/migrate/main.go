package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/config"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/database"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(context.Background(), appConfig.DatabaseURL); err != nil {
		log.Fatal(fmt.Errorf("run database migrations: %w", err))
	}
}

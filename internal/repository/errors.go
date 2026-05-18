package repository

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

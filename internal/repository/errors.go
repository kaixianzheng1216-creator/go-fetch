package repository

import (
	"errors"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"

	"github.com/jackc/pgx/v5"
)

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

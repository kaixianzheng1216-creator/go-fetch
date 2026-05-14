package store

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("未找到")

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

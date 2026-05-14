package server

import (
	"errors"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"

	"github.com/danielgtaylor/huma/v2"
)

func isStoreNotFound(err error) bool {
	return errors.Is(err, store.ErrNotFound)
}

func websiteLookupError(err error) error {
	if isStoreNotFound(err) {
		return huma.Error404NotFound("website not found")
	}

	return huma.Error500InternalServerError("failed to load website")
}

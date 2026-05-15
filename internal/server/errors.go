package server

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
)

func isNotFound(err error) bool {
	return errors.Is(err, store.ErrNotFound)
}

func websiteLookupError(err error) error {
	if isNotFound(err) {
		return huma.Error404NotFound("website not found")
	}

	return huma.Error500InternalServerError("load website failed")
}

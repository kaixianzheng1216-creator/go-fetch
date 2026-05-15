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
		return huma.Error404NotFound("站点不存在")
	}

	return huma.Error500InternalServerError("加载站点失败")
}

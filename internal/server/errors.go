package server

import (
	"errors"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"

	"github.com/danielgtaylor/huma/v2"
)

func isNotFound(err error) bool {
	return errors.Is(err, store.ErrNotFound)
}

func websiteLookupError(err error) error {
	if isNotFound(err) {
		return huma.Error404NotFound("网站不存在")
	}

	return huma.Error500InternalServerError("加载网站失败")
}

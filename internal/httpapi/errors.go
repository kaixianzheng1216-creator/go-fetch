package httpapi

import (
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

func isNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}

func parseUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, huma.Error400BadRequest(fmt.Sprintf("%s 必须是有效的 UUID", field))
	}
	return id, nil
}

func websiteLookupError(err error) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound("站点不存在")
	}
	return huma.Error500InternalServerError("加载站点失败")
}

func statsError(err error, fallbackMessage string) error {
	if errors.Is(err, domain.ErrUnsupportedMetricType) || errors.Is(err, service.ErrInvalidDateRange) {
		return huma.Error400BadRequest(err.Error())
	}
	if service.IsWebsiteAccessError(err) {
		return websiteLookupError(err)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}

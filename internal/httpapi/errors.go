package httpapi

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

func isNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
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

func websiteMutationError(err error, fallbackMessage string) error {
	if errors.Is(err, service.ErrInvalidWebsiteName) {
		return huma.Error400BadRequest("站点名称不能为空")
	}
	return websiteLookupErrorWithFallback(err, fallbackMessage)
}

func websiteLookupErrorWithFallback(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound("站点不存在")
	}
	return huma.Error500InternalServerError(fallbackMessage)
}

func collectionError(err error) error {
	switch {
	case errors.Is(err, service.ErrUnsupportedCollectionType):
		return huma.Error400BadRequest("不支持的采集类型")
	case errors.Is(err, service.ErrMissingClientInfo):
		return huma.Error500InternalServerError("读取请求失败")
	case isNotFound(err):
		return huma.Error400BadRequest("站点不存在")
	default:
		return huma.Error500InternalServerError("保存事件失败")
	}
}

func loginError(err error) error {
	if errors.Is(err, service.ErrInvalidCredentials) {
		return huma.Error401Unauthorized("用户名或密码错误")
	}
	return huma.Error500InternalServerError("加载用户失败")
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

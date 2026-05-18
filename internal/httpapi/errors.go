package httpapi

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

const (
	errorMessageAPIEndpointNotFound    = "接口不存在"
	errorMessageCurrentUserLoadFailed  = "加载当前用户失败"
	errorMessageDashboardBuildMissing  = "dashboard build is missing"
	errorMessageEventSaveFailed        = "保存事件失败"
	errorMessageInvalidCredentials     = "用户名或密码错误"
	errorMessageLoginSessionCreate     = "创建登录会话失败"
	errorMessageLogoutFailed           = "退出登录失败"
	errorMessageMetricsLoadFailed      = "加载指标数据失败"
	errorMessagePageviewsLoadFailed    = "加载页面浏览量失败"
	errorMessageRequestReadFailed      = "读取请求失败"
	errorMessageStatsLoadFailed        = "加载统计数据失败"
	errorMessageTrackerScriptMissing   = "tracking script is missing"
	errorMessageUnauthenticated        = "未登录"
	errorMessageUnsupportedCollectType = "不支持的采集类型"
	errorMessageUserLoadFailed         = "加载用户失败"
	errorMessageWebsiteListLoadFailed  = "加载站点列表失败"
	errorMessageWebsiteLoadFailed      = "加载站点失败"
	errorMessageWebsiteNameCannotEmpty = "站点名称不能为空"
	errorMessageWebsiteNotFound        = "站点不存在"
	errorMessageWebsiteCreateFailed    = "创建站点失败"
	errorMessageWebsiteUpdateFailed    = "更新站点失败"
)

func isNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}

func websiteLookupError(err error) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound(errorMessageWebsiteNotFound)
	}
	return huma.Error500InternalServerError(errorMessageWebsiteLoadFailed)
}

func websiteMutationError(err error, fallbackMessage string) error {
	if errors.Is(err, service.ErrInvalidWebsiteName) {
		return huma.Error400BadRequest(errorMessageWebsiteNameCannotEmpty)
	}
	return websiteLookupErrorWithFallback(err, fallbackMessage)
}

func websiteLookupErrorWithFallback(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound(errorMessageWebsiteNotFound)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}

func collectionError(err error) error {
	switch {
	case errors.Is(err, service.ErrUnsupportedCollectionType):
		return huma.Error400BadRequest(errorMessageUnsupportedCollectType)
	case errors.Is(err, service.ErrMissingClientInfo):
		return huma.Error500InternalServerError(errorMessageRequestReadFailed)
	case isNotFound(err):
		return huma.Error400BadRequest(errorMessageWebsiteNotFound)
	default:
		return huma.Error500InternalServerError(errorMessageEventSaveFailed)
	}
}

func loginError(err error) error {
	if errors.Is(err, service.ErrInvalidCredentials) {
		return huma.Error401Unauthorized(errorMessageInvalidCredentials)
	}
	return huma.Error500InternalServerError(errorMessageUserLoadFailed)
}

func statsError(err error, fallbackMessage string) error {
	if errors.Is(err, domain.ErrUnsupportedDateUnit) || errors.Is(err, domain.ErrUnsupportedMetricType) || errors.Is(err, service.ErrInvalidDateRange) {
		return huma.Error400BadRequest(err.Error())
	}
	if service.IsWebsiteAccessError(err) {
		return websiteLookupError(err)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}

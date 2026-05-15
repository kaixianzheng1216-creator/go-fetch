package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
)

func init() {
	huma.NewError = newLocalizedHumaError
}

func isNotFound(err error) bool {
	return errors.Is(err, store.ErrNotFound)
}

func websiteLookupError(err error) error {
	if isNotFound(err) {
		return huma.Error404NotFound("站点不存在")
	}

	return huma.Error500InternalServerError("加载站点失败")
}

func newLocalizedHumaError(status int, msg string, errs ...error) huma.StatusError {
	details := make([]*huma.ErrorDetail, len(errs))
	for i := range errs {
		if converted, ok := errs[i].(huma.ErrorDetailer); ok {
			details[i] = localizedHumaErrorDetail(converted.ErrorDetail())
			continue
		}
		if errs[i] != nil {
			details[i] = &huma.ErrorDetail{Message: localizedValidationMessage(errs[i].Error())}
		}
	}

	return &huma.ErrorModel{
		Status: status,
		Title:  localizedStatusText(status),
		Detail: localizedErrorMessage(msg),
		Errors: details,
	}
}

func localizedStatusText(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "请求错误"
	case http.StatusUnauthorized:
		return "未授权"
	case http.StatusForbidden:
		return "禁止访问"
	case http.StatusNotFound:
		return "未找到"
	case http.StatusMethodNotAllowed:
		return "请求方法不允许"
	case http.StatusRequestTimeout:
		return "请求超时"
	case http.StatusRequestEntityTooLarge:
		return "请求体过大"
	case http.StatusUnprocessableEntity:
		return "请求参数无法处理"
	case http.StatusInternalServerError:
		return "服务器内部错误"
	default:
		return fmt.Sprintf("HTTP %d 错误", status)
	}
}

func localizedErrorMessage(msg string) string {
	if msg == "validation failed" {
		return "请求参数校验失败"
	}
	return msg
}

func localizedHumaErrorDetail(detail *huma.ErrorDetail) *huma.ErrorDetail {
	if detail == nil {
		return nil
	}

	localized := *detail
	localized.Message = localizedValidationMessage(detail.Message)
	return &localized
}

func localizedValidationMessage(message string) string {
	switch {
	case message == "invalid integer":
		return "整数格式不正确"
	case message == "invalid float":
		return "小数格式不正确"
	case message == "invalid boolean":
		return "布尔值格式不正确"
	case strings.HasPrefix(message, "invalid date/time"):
		return "日期时间格式不正确"
	case message == "invalid url.URL value":
		return "URL 格式不正确"
	case strings.HasPrefix(message, "invalid value: invalid UUID length"):
		return "UUID 长度不正确"
	case strings.HasPrefix(message, "invalid value:"):
		return "值不合法"
	case message == "required query parameter is missing":
		return "缺少必填查询参数"
	case message == "required header parameter is missing":
		return "缺少必填请求头参数"
	case message == "required path parameter is missing":
		return "缺少必填路径参数"
	case strings.HasPrefix(message, "expected required property"):
		return "缺少必填字段"
	case strings.HasPrefix(message, "expected length <="):
		return "长度超过最大限制"
	case strings.HasPrefix(message, "expected length >="):
		return "长度小于最小限制"
	case strings.HasPrefix(message, "expected number <="):
		return "数值超过最大限制"
	case strings.HasPrefix(message, "expected number >="):
		return "数值小于最小限制"
	case strings.HasPrefix(message, "expected at most one value"):
		return "最多只能提供一个值"
	case strings.HasPrefix(message, "unknown query parameter"):
		return "未知查询参数"
	case message == "File required":
		return "文件不能为空"
	default:
		return message
	}
}

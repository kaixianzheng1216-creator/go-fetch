package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type CollectHandler struct {
	collect            service.Collect
	requestFromContext func(context.Context) *http.Request
	isNotFound         func(error) bool
}

func NewCollect(collect service.Collect, requestFromContext func(context.Context) *http.Request, isNotFound func(error) bool) CollectHandler {
	return CollectHandler{collect: collect, requestFromContext: requestFromContext, isNotFound: isNotFound}
}

type collectRequest struct {
	Body CollectRequest
}

func (handler CollectHandler) CollectEvent(ctx context.Context, input *collectRequest) (*OKOutput, error) {
	err := handler.collect.Collect(ctx, handler.requestFromContext(ctx), string(input.Body.Type), ToCollectPayload(input.Body.Payload))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnsupportedCollectionType):
			return nil, huma.Error400BadRequest("不支持的采集类型")
		case errors.Is(err, service.ErrMissingHTTPRequest):
			return nil, huma.Error500InternalServerError("读取请求失败")
		case handler.isNotFound(err):
			return nil, huma.Error400BadRequest("站点不存在")
		default:
			return nil, huma.Error500InternalServerError("保存事件失败")
		}
	}

	return NewOKOutput(), nil
}

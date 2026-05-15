package events

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/collector"
	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/event"
	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/website"
)

type Store interface {
	GetWebsiteForCollection(ctx context.Context, websiteID string) (websitedomain.Website, error)
	SaveEvent(ctx context.Context, event eventdomain.EventInput) error
}

type Handler struct {
	store              Store
	requestFromContext func(context.Context) *http.Request
	isNotFound         func(error) bool
}

func New(
	dataStore Store,
	requestFromContext func(context.Context) *http.Request,
	isNotFound func(error) bool,
) Handler {
	return Handler{
		store:              dataStore,
		requestFromContext: requestFromContext,
		isNotFound:         isNotFound,
	}
}

type collectRequest struct {
	Body CollectRequest
}

func (handler Handler) CollectEvent(ctx context.Context, request *collectRequest) (*okResponse, error) {
	collectionType, isSupportedCollectionType := eventdomain.ParseCollectionType(string(request.Body.Type))
	if !isSupportedCollectionType {
		return nil, huma.Error400BadRequest("不支持的采集类型")
	}

	request.Body.Type = CollectionType(collectionType)

	payload := ToCollectPayload(request.Body.Payload)
	if _, err := handler.store.GetWebsiteForCollection(ctx, payload.WebsiteID); err != nil {
		if handler.isNotFound(err) {
			return nil, huma.Error400BadRequest("站点不存在")
		}

		return nil, huma.Error500InternalServerError("加载站点失败")
	}

	httpRequest := handler.requestFromContext(ctx)
	if httpRequest == nil {
		return nil, huma.Error500InternalServerError("读取请求失败")
	}

	if collector.IsBot(httpRequest.UserAgent()) {
		return newOKResponse(), nil
	}

	if err := handler.store.SaveEvent(ctx, collector.BuildEventInput(httpRequest, payload, time.Now())); err != nil {
		return nil, huma.Error500InternalServerError("保存事件失败")
	}

	return newOKResponse(), nil
}

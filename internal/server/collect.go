package server

import (
	"context"
	"net/http"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/collector"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"

	"github.com/danielgtaylor/huma/v2"
)

func registerCollectRoutes(api huma.API, app *App) {
	collectOp := operation(
		http.MethodPost,
		"/api/collect",
		"collect",
		"Collection",
		http.StatusBadRequest,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	collectOp.MaxBodyBytes = 256 * 1024

	huma.Register(api, collectOp, app.collect)
}

func (a *App) collect(ctx context.Context, input *collectInput) (*jsonBody[httpapi.OK], error) {
	collectionType, ok := domain.ParseCollectionType(string(input.Body.Type))
	if !ok {
		return nil, huma.Error400BadRequest("不支持的采集类型")
	}

	input.Body.Type = httpapi.CollectionType(collectionType)

	payload := httpapi.CollectPayloadToDomain(input.Body.Payload)
	if _, err := a.store.GetWebsiteForCollection(ctx, payload.WebsiteID); err != nil {
		if isStoreNotFound(err) {
			return nil, huma.Error400BadRequest("网站不存在")
		}

		return nil, huma.Error500InternalServerError("加载网站失败")
	}

	r := requestFromContext(ctx)
	if r == nil {
		return nil, huma.Error500InternalServerError("读取请求失败")
	}

	if collector.IsBot(r.UserAgent()) {
		return &jsonBody[httpapi.OK]{Body: httpapi.OK{OK: true}}, nil
	}

	if err := a.store.SaveEvent(ctx, collector.BuildEventInput(r, payload, time.Now())); err != nil {
		return nil, huma.Error500InternalServerError("保存事件失败")
	}

	return &jsonBody[httpapi.OK]{Body: httpapi.OK{OK: true}}, nil
}

package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

const maxCollectBodyBytes = 256 * 1024

func registerCollectRoutes(api huma.API, collectHandler handler.CollectHandler) {
	op := huma.Operation{
		Method:      http.MethodPost,
		Path:        "/api/collect",
		OperationID: "collect",
		Summary:     "采集事件",
		Tags:        []string{"Collection"},
	}
	op.MaxBodyBytes = maxCollectBodyBytes
	huma.Register(api, op, collectHandler.CollectEvent)
}

package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

const maxCollectBodyBytes = 256 * 1024

func registerCollectRoutes(api huma.API, collectHandler handler.CollectHandler) {
	op := operation(
		http.MethodPost,
		"/api/collect",
		"collect",
		"采集事件",
		"Collection",
	)
	op.MaxBodyBytes = maxCollectBodyBytes
	huma.Register(api, op, collectHandler.CollectEvent)
}

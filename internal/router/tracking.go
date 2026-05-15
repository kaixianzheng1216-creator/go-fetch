package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

func RegisterCollect(api huma.API, collectHandler handler.CollectHandler) {
	huma.Register(
		api,
		WithMaxBodyBytes(
			NewOperation(
				http.MethodPost,
				"/api/collect",
				"collect",
				"采集事件",
				"Collection",
				http.StatusBadRequest,
				http.StatusUnprocessableEntity,
				http.StatusInternalServerError,
			),
			256*1024,
		),
		collectHandler.CollectEvent,
	)
}

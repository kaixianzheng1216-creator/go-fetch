package events

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"
)

func Register(api huma.API, handler Handler) {
	collectOp := httpapi.NewOperation(
		http.MethodPost,
		"/api/collect",
		"collect",
		"Collection",
		http.StatusBadRequest,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	collectOp.MaxBodyBytes = 256 * 1024

	huma.Register(api, collectOp, handler.CollectEvent)
}

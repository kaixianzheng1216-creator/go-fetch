package server

import (
	"net/http"
	"strings"

	assets "github.com/kaixianzheng1216-creator/go-fetch/internal/static"
)

const (
	contentTypeHTML = "text/html; charset=utf-8"
	contentTypeJS   = "application/javascript; charset=utf-8"

	apiPrefix = "/api/"
)

func (app *App) handleFrontendAsset(responseWriter http.ResponseWriter, request *http.Request) {
	http.FileServer(http.FS(assets.DistFS())).ServeHTTP(responseWriter, request)
}

func (app *App) handleScript(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", contentTypeJS)

	http.ServeFileFS(responseWriter, request, assets.StaticFS(), "script.js")
}

func (app *App) handleFrontend(responseWriter http.ResponseWriter, request *http.Request) {
	switch {
	case request.Method != http.MethodGet:
		http.NotFound(responseWriter, request)
		return

	case strings.HasPrefix(request.URL.Path, apiPrefix):
		writeProblemError(responseWriter, http.StatusNotFound, "接口不存在")
		return
	}

	indexHTML, err := assets.IndexHTML()
	if err != nil {
		http.Error(responseWriter, "前端构建产物不存在", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", contentTypeHTML)

	_, _ = responseWriter.Write(indexHTML)
}

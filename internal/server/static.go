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

func (a *App) handleFrontendAsset(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(assets.DistFS())).ServeHTTP(w, r)
}

func (a *App) handleScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentTypeJS)

	http.ServeFileFS(w, r, assets.StaticFS(), "script.js")
}

func (a *App) handleFrontend(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method != http.MethodGet:
		http.NotFound(w, r)
		return

	case strings.HasPrefix(r.URL.Path, apiPrefix):
		writeProblemError(w, http.StatusNotFound, "接口不存在")
		return
	}

	indexHTML, err := assets.IndexHTML()
	if err != nil {
		http.Error(w, "前端构建产物不存在", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentTypeHTML)

	_, _ = w.Write(indexHTML)
}

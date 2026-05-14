package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/web"
)

const (
	contentTypeHTML        = "text/html; charset=utf-8"
	contentTypeJS          = "application/javascript; charset=utf-8"
	contentTypeProblemJSON = "application/problem+json"

	apiPrefix = "/api/"
)

func (a *App) handleFrontendAsset(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(web.DistFS())).ServeHTTP(w, r)
}

func (a *App) handleScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentTypeJS)

	http.ServeFileFS(w, r, web.StaticFS(), "script.js")
}

func (a *App) handleFrontend(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method != http.MethodGet:
		http.NotFound(w, r)
		return

	case strings.HasPrefix(r.URL.Path, apiPrefix):
		writeProblemError(w, http.StatusNotFound, "not found")
		return
	}

	indexHTML, err := web.IndexHTML()
	if err != nil {
		http.Error(w, "frontend build not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentTypeHTML)

	_, _ = w.Write(indexHTML)
}

func writeProblemError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", contentTypeProblemJSON)
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(huma.ErrorModel{
		Title:  http.StatusText(status),
		Status: status,
		Detail: message,
	})
}

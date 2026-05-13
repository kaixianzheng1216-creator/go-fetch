package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/web"

	"github.com/danielgtaylor/huma/v2"
)

func (a *App) handleFrontendAsset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.FileServer(http.FS(web.DistFS())).ServeHTTP(w, r)
}

func (a *App) handleScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFileFS(w, r, web.StaticFS(), "script.js")
}

func (a *App) handleFrontend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api/") {
		writeProblemError(w, http.StatusNotFound, "not found")
		return
	}

	html, err := web.IndexHTML()
	if err != nil {
		http.Error(w, "frontend build not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(html)
}

func writeProblemError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(huma.ErrorModel{
		Title:  http.StatusText(status),
		Status: status,
		Detail: message,
	})
}

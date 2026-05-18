package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"

	webassets "github.com/kaixianzheng1216-creator/go-fetch/web"
)

func (srv server) registerAssets(chiRouter chi.Router) {
	chiRouter.Get("/assets/*", http.FileServer(http.FS(webassets.DashboardFS())).ServeHTTP)
	chiRouter.Get("/script.js", func(responseWriter http.ResponseWriter, _ *http.Request) {
		script, err := webassets.TrackerScript()
		if err != nil {
			http.Error(responseWriter, errorMessageTrackerScriptMissing, http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		if _, err := responseWriter.Write(script); err != nil {
			slog.Debug("write tracker script", "error", err)
		}
	})
	chiRouter.Get("/*", spaHandler)
}

func spaHandler(responseWriter http.ResponseWriter, request *http.Request) {
	switch {
	case request.Method != http.MethodGet:
		http.NotFound(responseWriter, request)
		return
	case strings.HasPrefix(request.URL.Path, "/api/"):
		responseWriter.Header().Set("Content-Type", "application/problem+json")
		responseWriter.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(responseWriter).Encode(huma.ErrorModel{
			Title:  errorMessageAPIEndpointNotFound,
			Status: http.StatusNotFound,
			Detail: errorMessageAPIEndpointNotFound,
		}); err != nil {
			slog.Debug("write API not found response", "error", err)
		}
		return
	}

	indexHTML, err := webassets.IndexHTML()
	if err != nil {
		http.Error(responseWriter, errorMessageDashboardBuildMissing, http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := responseWriter.Write(indexHTML); err != nil {
		slog.Debug("write dashboard index", "error", err)
	}
}

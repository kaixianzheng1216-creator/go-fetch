package server

import (
	"net/http"
	"strings"

	webassets "github.com/kaixianzheng1216-creator/go-fetch/web"
)

const (
	contentTypeHTML = "text/html; charset=utf-8"
	contentTypeJS   = "application/javascript; charset=utf-8"

	apiPrefix = "/api/"
)

func (server *Server) handleFrontendAsset(responseWriter http.ResponseWriter, request *http.Request) {
	http.FileServer(http.FS(webassets.DashboardFS())).ServeHTTP(responseWriter, request)
}

func (server *Server) handleScript(responseWriter http.ResponseWriter, _ *http.Request) {
	script, err := webassets.TrackerScript()
	if err != nil {
		http.Error(responseWriter, "tracking script is missing", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", contentTypeJS)
	_, _ = responseWriter.Write(script)
}

func (server *Server) handleFrontend(responseWriter http.ResponseWriter, request *http.Request) {
	switch {
	case request.Method != http.MethodGet:
		http.NotFound(responseWriter, request)
		return
	case strings.HasPrefix(request.URL.Path, apiPrefix):
		writeProblemError(responseWriter, http.StatusNotFound, "接口不存在")
		return
	}

	indexHTML, err := webassets.IndexHTML()
	if err != nil {
		http.Error(responseWriter, "dashboard build is missing", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", contentTypeHTML)
	_, _ = responseWriter.Write(indexHTML)
}

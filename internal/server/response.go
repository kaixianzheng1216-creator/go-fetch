package server

import (
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

const contentTypeProblemJSON = "application/problem+json"

func writeProblemError(responseWriter http.ResponseWriter, status int, message string) {
	responseWriter.Header().Set("Content-Type", contentTypeProblemJSON)
	responseWriter.WriteHeader(status)

	_ = json.NewEncoder(responseWriter).Encode(huma.ErrorModel{
		Title:  message,
		Status: status,
		Detail: message,
	})
}

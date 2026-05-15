package server

import (
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

const contentTypeProblemJSON = "application/problem+json"

func writeProblemError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", contentTypeProblemJSON)
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(huma.ErrorModel{
		Title:  message,
		Status: status,
		Detail: message,
	})
}

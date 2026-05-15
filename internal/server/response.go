package server

import (
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

const contentTypeProblemJSON = "application/problem+json"

type OK struct {
	OK bool `json:"ok"`
}

type jsonBody[T any] struct {
	Body T
}

func jsonResponse[T any](body T) *jsonBody[T] {
	return &jsonBody[T]{Body: body}
}

func okResponse() *jsonBody[OK] {
	response := OK{OK: true}

	return jsonResponse(response)
}

func writeProblemError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", contentTypeProblemJSON)
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(huma.ErrorModel{
		Title:  message,
		Status: status,
		Detail: message,
	})
}

package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"go-fetch/internal/httpapi"
	"go-fetch/internal/store"

	"github.com/danielgtaylor/huma/v2"
)

func init() {
	huma.NewError = newAPIError
	huma.NewErrorWithContext = func(_ huma.Context, status int, msg string, errs ...error) huma.StatusError {
		return newAPIError(status, msg, errs...)
	}
}

type apiError struct {
	ErrorDetail httpapi.ErrorDetail `json:"error"`
	status      int
}

func newAPIError(status int, msg string, errs ...error) huma.StatusError {
	if msg == "" {
		msg = http.StatusText(status)
	}
	if msg == "validation failed" {
		for _, err := range errs {
			if err == nil {
				continue
			}
			text := err.Error()
			if strings.Contains(text, "unexpected EOF") ||
				strings.Contains(text, "unexpected end") ||
				strings.Contains(text, "invalid character") ||
				strings.Contains(text, "cannot unmarshal") {
				status = http.StatusBadRequest
				msg = "invalid json"
				break
			}
		}
	}
	return &apiError{
		status:      status,
		ErrorDetail: httpapi.ErrorDetail{Message: msg},
	}
}

func (e *apiError) Error() string {
	return e.ErrorDetail.Message
}

func (e *apiError) GetStatus() int {
	return e.status
}

func (e *apiError) ContentType(string) string {
	return "application/json"
}

func (e *apiError) MarshalJSON() ([]byte, error) {
	return json.Marshal(httpapi.ErrorResponse{Error: e.ErrorDetail})
}

func (e *apiError) Schema(registry huma.Registry) *huma.Schema {
	return huma.SchemaFromType(registry, reflect.TypeOf(httpapi.ErrorResponse{}))
}

func isStoreNotFound(err error) bool {
	return errors.Is(err, store.ErrNotFound)
}

func websiteNotFound() error {
	return huma.Error404NotFound("website not found")
}

func websiteLookupError(err error) error {
	if isStoreNotFound(err) {
		return websiteNotFound()
	}
	return huma.Error500InternalServerError("failed to load website")
}

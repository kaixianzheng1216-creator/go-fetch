package httpapi

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

const (
	errorMessageAPIEndpointNotFound    = "API endpoint not found"
	errorMessageCurrentUserLoadFailed  = "failed to load current user"
	errorMessageDashboardBuildMissing  = "dashboard build is missing"
	errorMessageEventSaveFailed        = "failed to save event"
	errorMessageInvalidCredentials     = "invalid username or password"
	errorMessageLoginSessionCreate     = "failed to create login session"
	errorMessageLogoutFailed           = "failed to log out"
	errorMessageMetricsLoadFailed      = "failed to load metrics"
	errorMessagePageviewsLoadFailed    = "failed to load pageviews"
	errorMessageRequestReadFailed      = "failed to read request"
	errorMessageStatsLoadFailed        = "failed to load stats"
	errorMessageTrackerScriptMissing   = "tracking script is missing"
	errorMessageUnauthenticated        = "authentication required"
	errorMessageUnsupportedEventType   = "unsupported event type"
	errorMessageUserLoadFailed         = "failed to load user"
	errorMessageWebsiteListLoadFailed  = "failed to load websites"
	errorMessageWebsiteLoadFailed      = "failed to load website"
	errorMessageWebsiteNameCannotEmpty = "website name cannot be empty"
	errorMessageWebsiteNotFound        = "website not found"
	errorMessageWebsiteCreateFailed    = "failed to create website"
	errorMessageWebsiteUpdateFailed    = "failed to update website"
)

func isNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}

func websiteLookupError(err error) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound(errorMessageWebsiteNotFound)
	}
	return huma.Error500InternalServerError(errorMessageWebsiteLoadFailed)
}

func websiteMutationError(err error, fallbackMessage string) error {
	if errors.Is(err, service.ErrInvalidWebsiteName) {
		return huma.Error400BadRequest(errorMessageWebsiteNameCannotEmpty)
	}
	return websiteLookupErrorWithFallback(err, fallbackMessage)
}

func websiteLookupErrorWithFallback(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound(errorMessageWebsiteNotFound)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}

func collectionError(err error) error {
	switch {
	case errors.Is(err, service.ErrUnsupportedEventType):
		return huma.Error400BadRequest(errorMessageUnsupportedEventType)
	case errors.Is(err, service.ErrMissingClientInfo):
		return huma.Error500InternalServerError(errorMessageRequestReadFailed)
	case isNotFound(err):
		return huma.Error400BadRequest(errorMessageWebsiteNotFound)
	default:
		return huma.Error500InternalServerError(errorMessageEventSaveFailed)
	}
}

func loginError(err error) error {
	if errors.Is(err, service.ErrInvalidCredentials) {
		return huma.Error401Unauthorized(errorMessageInvalidCredentials)
	}
	return huma.Error500InternalServerError(errorMessageUserLoadFailed)
}

func statsError(err error, fallbackMessage string) error {
	if errors.Is(err, domain.ErrUnsupportedDateUnit) || errors.Is(err, domain.ErrUnsupportedMetricType) || errors.Is(err, service.ErrInvalidDateRange) {
		return huma.Error400BadRequest(err.Error())
	}
	if service.IsWebsiteAccessError(err) {
		return websiteLookupError(err)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}

package httpapi

import (
	"errors"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
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

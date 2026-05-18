package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/netip"
	"net/url"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

const maxCollectBodyBytes = 256 * 1024

type collectInput struct {
	Body struct {
		Type       collectEventTypeParam `json:"type" required:"true"`
		WebsiteID  uuid.UUID             `json:"websiteId" required:"true" format:"uuid"`
		URL        string                `json:"url" required:"true" minLength:"1"`
		Referrer   string                `json:"referrer,omitempty"`
		Title      string                `json:"title,omitempty"`
		Screen     string                `json:"screen,omitempty"`
		Language   string                `json:"language,omitempty"`
		DistinctID string                `json:"distinctId,omitempty" maxLength:"50"`
		Name       string                `json:"name,omitempty"`
		Data       map[string]any        `json:"data,omitempty"`
	}
}

type collectEventTypeParam string

func (collectEventTypeParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(domain.TrackedEventTypeValues()),
	}
}

func (srv server) registerCollectRoutes(humaAPI huma.API) {
	operation := publicOperation(http.MethodPost, "/api/collect", "collect", "Collect event", "Collection")
	operation.MaxBodyBytes = maxCollectBodyBytes
	huma.Register(humaAPI, operation, srv.collectEvent)
}

func (srv server) collectEvent(ctx context.Context, input *collectInput) (*okOutput, error) {
	err := srv.collection.CollectEvent(ctx, service.CollectEventInput{
		Client: srv.clientInfoFromRequest(requestFromContext(ctx)),
		Event: domain.TrackedEvent{
			Type:       domain.TrackedEventType(input.Body.Type),
			WebsiteID:  input.Body.WebsiteID,
			URL:        input.Body.URL,
			Referrer:   input.Body.Referrer,
			Title:      input.Body.Title,
			Screen:     input.Body.Screen,
			Language:   input.Body.Language,
			DistinctID: input.Body.DistinctID,
			Name:       input.Body.Name,
			Data:       input.Body.Data,
		},
	})
	if err != nil {
		return nil, collectionError(err)
	}

	return newOKOutput(), nil
}

func (srv server) clientInfoFromRequest(request *http.Request) service.ClientInfo {
	if request == nil {
		return service.ClientInfo{}
	}

	client := service.ClientInfo{
		IP:        clientIP(request.RemoteAddr),
		UserAgent: request.UserAgent(),
	}
	if srv.config.TrustProxyHeaders {
		client.Country = countryHeader(request.Header)
		client.Region = geoHeader(request.Header, "CF-IPRegionCode", "CF-IPRegion", "X-Vercel-IP-Country-Region", "X-Appengine-Region", "CloudFront-Viewer-Country-Region", "X-Geo-Region")
		client.City = geoHeader(request.Header, "CF-IPCity", "X-Vercel-IP-City", "X-Appengine-City", "CloudFront-Viewer-City", "X-Geo-City")
	}

	return client
}

func clientIP(remoteAddr string) string {
	addrPort, err := netip.ParseAddrPort(remoteAddr)
	if err == nil {
		return addrPort.Addr().String()
	}

	addr, err := netip.ParseAddr(remoteAddr)
	if err == nil {
		return addr.String()
	}

	return remoteAddr
}

func countryHeader(header http.Header) string {
	country := strings.ToUpper(geoHeader(header, "CF-IPCountry", "X-Vercel-IP-Country", "X-Appengine-Country", "CloudFront-Viewer-Country", "X-Country-Code", "X-Geo-Country"))
	if country == "XX" || len(country) != 2 || !isASCIILetters(country) {
		return ""
	}
	return country
}

func geoHeader(header http.Header, names ...string) string {
	for _, name := range names {
		value := cleanGeoHeaderValue(header.Get(name))
		if value != "" {
			return value
		}
	}
	return ""
}

func cleanGeoHeaderValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	unescaped, err := url.QueryUnescape(value)
	if err == nil {
		value = unescaped
	}
	return strings.TrimSpace(value)
}

func isASCIILetters(value string) bool {
	for _, r := range value {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
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

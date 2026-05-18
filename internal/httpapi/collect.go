package httpapi

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

const maxCollectBodyBytes = 256 * 1024

type collectInput struct {
	Body CollectEventRequest
}

type CollectEventRequest struct {
	Type    CollectionTypeParam        `json:"type,omitempty"`
	Payload CollectEventPayloadRequest `json:"payload" required:"true"`
}

type CollectionTypeParam string

func (CollectionTypeParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(domain.CollectionTypeValues()),
	}
}

type CollectEventPayloadRequest struct {
	WebsiteID  uuid.UUID      `json:"website" required:"true" format:"uuid"`
	URL        string         `json:"url" required:"true" minLength:"1"`
	Referrer   string         `json:"referrer,omitempty"`
	Title      string         `json:"title,omitempty"`
	Screen     string         `json:"screen,omitempty"`
	Language   string         `json:"language,omitempty"`
	DistinctID string         `json:"distinctId,omitempty" maxLength:"50"`
	Name       string         `json:"name,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

func (apiServer server) registerCollectRoutes(humaAPI huma.API) {
	operation := publicOperation(http.MethodPost, "/api/collect", "collect", "Collect event", "Collection")
	operation.MaxBodyBytes = maxCollectBodyBytes
	huma.Register(humaAPI, operation, apiServer.collectEvent)
}

func (apiServer server) collectEvent(ctx context.Context, input *collectInput) (*okOutput, error) {
	err := apiServer.collect.CollectEvent(ctx, service.CollectEventParams{
		Client:  apiServer.clientInfoFromRequest(requestFromContext(ctx)),
		Type:    domain.CollectionType(input.Body.Type),
		Payload: toDomainCollectPayload(input.Body.Payload),
	})
	if err != nil {
		return nil, collectionError(err)
	}

	return toOKOutput(), nil
}

func toDomainCollectPayload(payload CollectEventPayloadRequest) domain.CollectPayload {
	return domain.CollectPayload{
		WebsiteID:  payload.WebsiteID,
		URL:        payload.URL,
		Referrer:   payload.Referrer,
		Title:      payload.Title,
		Screen:     payload.Screen,
		Language:   payload.Language,
		DistinctID: payload.DistinctID,
		Name:       payload.Name,
		Data:       payload.Data,
	}
}

func (apiServer server) clientInfoFromRequest(request *http.Request) service.ClientInfo {
	if request == nil {
		return service.ClientInfo{}
	}

	client := service.ClientInfo{
		IP:        clientIP(request.RemoteAddr),
		UserAgent: request.UserAgent(),
	}
	if apiServer.config.TrustProxyHeaders {
		client.Country = countryHeader(request.Header)
		client.Region = geoHeader(request.Header, "CF-IPRegionCode", "CF-IPRegion", "X-Vercel-IP-Country-Region", "X-Appengine-Region", "CloudFront-Viewer-Country-Region", "X-Geo-Region")
		client.City = geoHeader(request.Header, "CF-IPCity", "X-Vercel-IP-City", "X-Appengine-City", "CloudFront-Viewer-City", "X-Geo-City")
	}

	return client
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return host
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

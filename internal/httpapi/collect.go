package httpapi

import (
	"context"
	"net"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

const maxCollectBodyBytes = 256 * 1024

type collectInput struct {
	Body CollectRequest
}

type CollectRequest struct {
	Type    CollectionTypeParam   `json:"type,omitempty"`
	Payload CollectPayloadRequest `json:"payload" required:"true"`
}

type CollectionTypeParam string

func (CollectionTypeParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(domain.CollectionTypeValues()),
	}
}

type CollectPayloadRequest struct {
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
	operation := publicOperation(http.MethodPost, "/api/collect", "collect", "采集事件", "Collection")
	operation.MaxBodyBytes = maxCollectBodyBytes
	huma.Register(humaAPI, operation, apiServer.collectEvent)
}

func (apiServer server) collectEvent(ctx context.Context, input *collectInput) (*okOutput, error) {
	err := apiServer.collect.Collect(ctx, service.CollectionParams{
		Client:  clientInfoFromRequest(requestFromContext(ctx)),
		Type:    domain.CollectionType(input.Body.Type),
		Payload: toDomainCollectPayload(input.Body.Payload),
	})
	if err != nil {
		return nil, collectionError(err)
	}

	return toOKOutput(), nil
}

func toDomainCollectPayload(payload CollectPayloadRequest) domain.CollectPayload {
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

func clientInfoFromRequest(request *http.Request) service.ClientInfo {
	if request == nil {
		return service.ClientInfo{}
	}
	return service.ClientInfo{
		IP:        clientIP(request.RemoteAddr),
		UserAgent: request.UserAgent(),
	}
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return host
	}
	return remoteAddr
}

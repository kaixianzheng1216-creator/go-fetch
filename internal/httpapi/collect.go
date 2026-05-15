package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

const maxCollectBodyBytes = 256 * 1024

type collectRequest struct {
	Body CollectRequest
}

type CollectRequest struct {
	Type    CollectionType `json:"type,omitempty"`
	Payload CollectPayload `json:"payload" required:"true"`
}

type CollectionType string

func (CollectionType) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(domain.CollectionTypeValues()),
	}
}

type CollectPayload struct {
	WebsiteID  string         `json:"website" required:"true" format:"uuid"`
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
	operation := huma.Operation{
		Method:      http.MethodPost,
		Path:        "/api/collect",
		OperationID: "collect",
		Summary:     "采集事件",
		Tags:        []string{"Collection"},
	}
	operation.MaxBodyBytes = maxCollectBodyBytes
	huma.Register(humaAPI, operation, apiServer.collectEvent)
}

func (apiServer server) collectEvent(ctx context.Context, input *collectRequest) (*okOutput, error) {
	payload, err := toCollectPayload(input.Body.Payload)
	if err != nil {
		return nil, err
	}

	err = apiServer.collect.Collect(ctx, requestFromContext(ctx), string(input.Body.Type), payload)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnsupportedCollectionType):
			return nil, huma.Error400BadRequest("不支持的采集类型")
		case errors.Is(err, service.ErrMissingHTTPRequest):
			return nil, huma.Error500InternalServerError("读取请求失败")
		case isNotFound(err):
			return nil, huma.Error400BadRequest("站点不存在")
		default:
			return nil, huma.Error500InternalServerError("保存事件失败")
		}
	}

	return toOKOutput(), nil
}

func toCollectPayload(payload CollectPayload) (domain.CollectPayload, error) {
	websiteID, err := parseUUID(payload.WebsiteID, "payload.website")
	if err != nil {
		return domain.CollectPayload{}, err
	}

	return domain.CollectPayload{
		WebsiteID:  websiteID,
		URL:        payload.URL,
		Referrer:   payload.Referrer,
		Title:      payload.Title,
		Screen:     payload.Screen,
		Language:   payload.Language,
		DistinctID: payload.DistinctID,
		Name:       payload.Name,
		Data:       payload.Data,
	}, nil
}

package event

import (
	"github.com/danielgtaylor/huma/v2"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/event"
)

type CollectRequest struct {
	Type    CollectionType `json:"type,omitempty"`
	Payload CollectPayload `json:"payload" required:"true"`
}

type CollectionType string

func (CollectionType) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(eventdomain.CollectionTypeValues()),
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

func ToCollectPayload(payload CollectPayload) eventdomain.CollectPayload {
	return eventdomain.CollectPayload{
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

func enumValues(values []string) []any {
	result := make([]any, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}

package httpapi

import (
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"

	"github.com/danielgtaylor/huma/v2"
)

type User struct {
	ID          string     `json:"id" format:"uuid"`
	Username    string     `json:"username"`
	LogoURL     string     `json:"logoUrl,omitempty"`
	DisplayName string     `json:"displayName,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" required:"true" minLength:"1"`
	Password string `json:"password" required:"true" minLength:"1" writeOnly:"true"`
}

type LoginResponse struct {
	User User `json:"user"`
}

type Website struct {
	ID        string    `json:"id" format:"uuid"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"createdAt"`
}

type WebsiteRequest struct {
	Name   string `json:"name" required:"true" minLength:"1" maxLength:"100"`
	Domain string `json:"domain,omitempty" maxLength:"500"`
}

type WebsiteStats struct {
	Pageviews       int64 `json:"pageviews"`
	Visitors        int64 `json:"visitors"`
	Visits          int64 `json:"visits"`
	Bounces         int64 `json:"bounces"`
	TotalTime       int64 `json:"totalTime"`
	AvgVisitSeconds int64 `json:"avgVisitSeconds"`
}

type PageviewPoint struct {
	Time     time.Time `json:"time"`
	Label    string    `json:"label"`
	Views    int64     `json:"views"`
	Visitors int64     `json:"visitors"`
}

type MetricRow struct {
	Name     string `json:"name"`
	Views    int64  `json:"views"`
	Visitors int64  `json:"visitors"`
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

type OK struct {
	OK bool `json:"ok"`
}

func UserFromDomain(user domain.User) User {
	return User{
		ID:          user.ID,
		Username:    user.Username,
		LogoURL:     user.LogoURL,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}
}

func WebsiteFromDomain(website domain.Website) Website {
	return Website{
		ID:        website.ID,
		Name:      website.Name,
		Domain:    website.Domain,
		CreatedAt: website.CreatedAt,
	}
}

func WebsiteStatsFromDomain(stats domain.WebsiteStats) WebsiteStats {
	return WebsiteStats{
		Pageviews:       stats.Pageviews,
		Visitors:        stats.Visitors,
		Visits:          stats.Visits,
		Bounces:         stats.Bounces,
		TotalTime:       stats.TotalTime,
		AvgVisitSeconds: stats.AvgVisitSeconds,
	}
}

func PageviewPointFromDomain(point domain.PageviewPoint) PageviewPoint {
	return PageviewPoint{
		Time:     point.Time,
		Label:    point.Label,
		Views:    point.Views,
		Visitors: point.Visitors,
	}
}

func PageviewPointsFromDomain(points []domain.PageviewPoint) []PageviewPoint {
	result := make([]PageviewPoint, 0, len(points))
	for _, point := range points {
		result = append(result, PageviewPointFromDomain(point))
	}

	return result
}

func MetricRowFromDomain(row domain.MetricRow) MetricRow {
	return MetricRow{
		Name:     row.Name,
		Views:    row.Views,
		Visitors: row.Visitors,
	}
}

func MetricRowsFromDomain(rows []domain.MetricRow) []MetricRow {
	result := make([]MetricRow, 0, len(rows))
	for _, row := range rows {
		result = append(result, MetricRowFromDomain(row))
	}

	return result
}

func WebsitesFromDomain(websites []domain.Website) []Website {
	result := make([]Website, 0, len(websites))
	for _, website := range websites {
		result = append(result, WebsiteFromDomain(website))
	}

	return result
}

func CollectPayloadToDomain(payload CollectPayload) domain.CollectPayload {
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

func enumValues(values []string) []any {
	result := make([]any, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}

	return result
}

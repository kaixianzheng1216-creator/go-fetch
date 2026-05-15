package server

import "time"

type User struct {
	ID        string     `json:"id" format:"uuid"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
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

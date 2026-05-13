package domain

import "time"

const (
	EventTypePageView EventType = 1
	EventTypeCustom   EventType = 2
)

type EventType int

type User struct {
	ID           string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type Website struct {
	ID        string
	Name      string
	Domain    string
	CreatedAt time.Time
}

type WebsiteStats struct {
	Pageviews       int64
	Visitors        int64
	Visits          int64
	Bounces         int64
	TotalTime       int64
	AvgVisitSeconds int64
}

type PageviewPoint struct {
	Time     time.Time
	Label    string
	Views    int64
	Visitors int64
}

type MetricRow struct {
	Name     string
	Views    int64
	Visitors int64
}

type CollectPayload struct {
	WebsiteID string
	URL       string
	Referrer  string
	Title     string
	Screen    string
	Language  string
	Name      string
	Data      map[string]any
}

type CollectResult struct {
	SessionID string
	VisitID   string
}

type EventInput struct {
	WebsiteID      string
	SessionID      string
	VisitID        string
	EventType      EventType
	EventName      string
	URLPath        string
	URLQuery       string
	ReferrerPath   string
	ReferrerDomain string
	PageTitle      string
	Hostname       string
	UTMSource      string
	UTMMedium      string
	UTMCampaign    string
	UTMContent     string
	UTMTerm        string
	Browser        string
	OS             string
	Device         string
	Screen         string
	Language       string
	Country        string
	CreatedAt      time.Time
	Data           map[string]any
}

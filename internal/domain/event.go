package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType int

const (
	EventTypePageView EventType = 1
	EventTypeCustom   EventType = 2
)

type TrackedEventType string

const (
	TrackedEventTypePageView TrackedEventType = "pageview"
	TrackedEventTypeCustom   TrackedEventType = "event"
)

type EventDataType int

const (
	EventDataTypeString  EventDataType = 1
	EventDataTypeNumber  EventDataType = 2
	EventDataTypeBoolean EventDataType = 3
	EventDataTypeDate    EventDataType = 4
	EventDataTypeArray   EventDataType = 5
)

type TrackedEvent struct {
	Type       TrackedEventType
	WebsiteID  uuid.UUID
	URL        string
	Referrer   string
	Title      string
	Screen     string
	Language   string
	DistinctID string
	Name       string
	Data       map[string]any
}

type EventRecord struct {
	WebsiteID      uuid.UUID
	SessionID      uuid.UUID
	VisitID        uuid.UUID
	EventType      EventType
	EventName      string
	URLPath        string
	URLQuery       string
	ReferrerPath   string
	ReferrerQuery  string
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
	Region         string
	City           string
	DistinctID     string
	CreatedAt      time.Time
	Data           map[string]any
}

type WebsiteStats struct {
	Pageviews       int64
	Visitors        int64
	Visits          int64
	Bounces         int64
	TotalTime       int64
	AvgVisitSeconds int64
}

type PageviewBucket struct {
	Time     time.Time
	Label    string
	Views    int64
	Visitors int64
}

type Metric struct {
	Name     string
	Views    int64
	Visitors int64
}

func NormalizeTrackedEventType(eventType TrackedEventType, eventName string) (TrackedEventType, bool) {
	if eventType == "" {
		if eventName != "" {
			return TrackedEventTypeCustom, true
		}
		return TrackedEventTypePageView, true
	}

	switch eventType {
	case TrackedEventTypePageView, TrackedEventTypeCustom:
		return eventType, true
	default:
		return "", false
	}
}

func TrackedEventTypeValues() []string {
	return []string{string(TrackedEventTypePageView), string(TrackedEventTypeCustom)}
}

func (eventType TrackedEventType) EventType() EventType {
	if eventType == TrackedEventTypeCustom {
		return EventTypeCustom
	}

	return EventTypePageView
}

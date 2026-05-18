package domain

import (
	"slices"
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

var trackedEventTypes = [...]TrackedEventType{
	TrackedEventTypePageView,
	TrackedEventTypeCustom,
}

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

func NormalizeTrackedEventType(eventType TrackedEventType, eventName string) (TrackedEventType, bool) {
	if eventType == "" {
		if eventName != "" {
			return TrackedEventTypeCustom, true
		}
		return TrackedEventTypePageView, true
	}

	if slices.Contains(trackedEventTypes[:], eventType) {
		return eventType, true
	}

	return "", false
}

func TrackedEventTypeValues() []string {
	return stringEnumValues(trackedEventTypes[:])
}

func (eventType TrackedEventType) EventType() EventType {
	if eventType == TrackedEventTypeCustom {
		return EventTypeCustom
	}

	return EventTypePageView
}

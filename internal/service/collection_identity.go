package service

import (
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/util"
)

type trackingIdentity struct {
	distinctID string
	sessionID  uuid.UUID
	visitID    uuid.UUID
}

func newTrackingIdentity(websiteID uuid.UUID, distinctID string, client trackingClient, now time.Time) trackingIdentity {
	distinctID = util.TruncateRunes(distinctID, maxDistinctIDLength)
	sessionID := sessionIDFor(websiteID, distinctID, client.ip, client.userAgent, now)
	return trackingIdentity{
		distinctID: distinctID,
		sessionID:  sessionID,
		visitID:    visitIDFor(sessionID, now),
	}
}

func sessionIDFor(websiteID uuid.UUID, distinctID, clientIP, userAgent string, now time.Time) uuid.UUID {
	value := websiteID.String() + "|" + visitorIdentity(distinctID, clientIP, userAgent) + "|" + now.UTC().Format(sessionWindowFormat)
	return stableUUID(value)
}

func visitIDFor(sessionID uuid.UUID, now time.Time) uuid.UUID {
	return stableUUID(sessionID.String() + "|" + strconv.FormatInt(now.Unix()/visitWindowSeconds, 10))
}

func stableUUID(value string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(value))
}

func visitorIdentity(distinctID, clientIP, userAgent string) string {
	if distinctID != "" {
		return distinctID
	}
	return clientIP + "|" + userAgent
}

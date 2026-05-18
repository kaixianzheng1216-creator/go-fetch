package service

import (
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/textutil"
)

const (
	sessionWindowFormat = "2006-01"
	visitWindowSeconds  = 1800

	defaultURLScheme = "https"
	defaultURLPath   = "/"

	maxEventNameLength  = 50
	maxURLPartLength    = 500
	maxPageTitleLength  = 500
	maxHostnameLength   = 100
	maxUTMValueLength   = 255
	maxBrowserLength    = 20
	maxOSLength         = 20
	maxDeviceLength     = 20
	maxScreenLength     = 11
	maxLanguageLength   = 35
	maxCountryLength    = 2
	maxRegionLength     = 20
	maxCityLength       = 50
	maxDistinctIDLength = 50

	laptopMaxScreenWidth = 1280
)

func buildEventRecord(client eventClient, event domain.TrackedEvent, website domain.Website, now time.Time) domain.EventRecord {
	eventURLs := newEventURLs(event, website)
	urlFields := newEventURLFields(eventURLs)
	utmFields := newUTMFields(eventURLs.page.Query())
	identity := newEventIdentity(event.WebsiteID, event.DistinctID, client, now)

	return domain.EventRecord{
		WebsiteID:      event.WebsiteID,
		SessionID:      identity.sessionID,
		VisitID:        identity.visitID,
		EventType:      event.Type.EventType(),
		EventName:      textutil.TruncateRunes(event.Name, maxEventNameLength),
		URLPath:        urlFields.path,
		URLQuery:       urlFields.query,
		ReferrerPath:   urlFields.referrerPath,
		ReferrerQuery:  urlFields.referrerQuery,
		ReferrerDomain: urlFields.referrerDomain,
		PageTitle:      textutil.TruncateRunes(event.Title, maxPageTitleLength),
		Hostname:       urlFields.hostname,
		UTMSource:      utmFields.source,
		UTMMedium:      utmFields.medium,
		UTMCampaign:    utmFields.campaign,
		UTMContent:     utmFields.content,
		UTMTerm:        utmFields.term,
		Browser:        textutil.TruncateRunes(client.browser, maxBrowserLength),
		OS:             textutil.TruncateRunes(client.os, maxOSLength),
		Device:         textutil.TruncateRunes(client.device, maxDeviceLength),
		Screen:         textutil.TruncateRunes(event.Screen, maxScreenLength),
		Language:       textutil.TruncateRunes(event.Language, maxLanguageLength),
		Country:        textutil.TruncateRunes(client.country, maxCountryLength),
		Region:         textutil.TruncateRunes(client.region, maxRegionLength),
		City:           textutil.TruncateRunes(client.city, maxCityLength),
		DistinctID:     identity.distinctID,
		CreatedAt:      now,
		Data:           event.Data,
	}
}

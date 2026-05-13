package server

import "go-fetch/internal/httpapi"

type emptyInput struct{}

type loginInput struct {
	Body httpapi.LoginRequest
}

type collectInput struct {
	Body httpapi.CollectRequest
}

type websiteBodyInput struct {
	Body httpapi.WebsiteRequest
}

type websitePathInput struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
}

type updateWebsiteInput struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
	Body      httpapi.WebsiteRequest
}

type dateRangeInput struct {
	WebsiteID string               `path:"websiteID" format:"uuid"`
	StartAt   optionalParam[int64] `query:"startAt"`
	EndAt     optionalParam[int64] `query:"endAt"`
}

type pageviewsInput struct {
	WebsiteID string               `path:"websiteID" format:"uuid"`
	StartAt   optionalParam[int64] `query:"startAt"`
	EndAt     optionalParam[int64] `query:"endAt"`
	Unit      dateUnitParam        `query:"unit"`
}

type metricsInput struct {
	WebsiteID string               `path:"websiteID" format:"uuid"`
	StartAt   optionalParam[int64] `query:"startAt"`
	EndAt     optionalParam[int64] `query:"endAt"`
	Type      metricTypeParam      `query:"type" required:"true"`
	Limit     metricLimit          `query:"limit"`
}

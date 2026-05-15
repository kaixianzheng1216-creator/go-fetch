package server

type emptyInput struct{}

type loginInput struct {
	Body LoginRequest
}

type websiteInput struct {
	Body WebsiteRequest
}

type websiteIDInput struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
}

type updateWebsiteInput struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
	Body      WebsiteRequest
}

type collectInput struct {
	Body CollectRequest
}

type statsInput struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
	StartAt   int64  `query:"startAt"`
	EndAt     int64  `query:"endAt"`
}

type pageviewsInput struct {
	WebsiteID string        `path:"websiteID" format:"uuid"`
	StartAt   int64         `query:"startAt"`
	EndAt     int64         `query:"endAt"`
	Unit      dateUnitParam `query:"unit"`
}

type metricsInput struct {
	WebsiteID string          `path:"websiteID" format:"uuid"`
	StartAt   int64           `query:"startAt"`
	EndAt     int64           `query:"endAt"`
	Type      metricTypeParam `query:"type" required:"true"`
	Limit     metricLimit     `query:"limit"`
}

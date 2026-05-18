package domain

import "time"

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

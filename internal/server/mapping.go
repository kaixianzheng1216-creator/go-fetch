package server

import "github.com/kaixianzheng1216-creator/go-fetch/internal/domain"

func toUser(user domain.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func toWebsite(website domain.Website) Website {
	return Website{
		ID:        website.ID,
		Name:      website.Name,
		Domain:    website.Domain,
		CreatedAt: website.CreatedAt,
	}
}

func toWebsites(websites []domain.Website) []Website {
	result := make([]Website, 0, len(websites))
	for _, website := range websites {
		result = append(result, toWebsite(website))
	}

	return result
}

func toWebsiteStats(stats domain.WebsiteStats) WebsiteStats {
	return WebsiteStats{
		Pageviews:       stats.Pageviews,
		Visitors:        stats.Visitors,
		Visits:          stats.Visits,
		Bounces:         stats.Bounces,
		TotalTime:       stats.TotalTime,
		AvgVisitSeconds: stats.AvgVisitSeconds,
	}
}

func toPageviewPoint(point domain.PageviewPoint) PageviewPoint {
	return PageviewPoint{
		Time:     point.Time,
		Label:    point.Label,
		Views:    point.Views,
		Visitors: point.Visitors,
	}
}

func toPageviewPoints(points []domain.PageviewPoint) []PageviewPoint {
	result := make([]PageviewPoint, 0, len(points))
	for _, point := range points {
		result = append(result, toPageviewPoint(point))
	}

	return result
}

func toMetricRow(row domain.MetricRow) MetricRow {
	return MetricRow{
		Name:     row.Name,
		Views:    row.Views,
		Visitors: row.Visitors,
	}
}

func toMetricRows(rows []domain.MetricRow) []MetricRow {
	result := make([]MetricRow, 0, len(rows))
	for _, row := range rows {
		result = append(result, toMetricRow(row))
	}

	return result
}

func toCollectPayload(payload CollectPayload) domain.CollectPayload {
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

package httpapi

import (
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

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

func WebsiteFromDomain(website domain.Website) Website {
	return Website{
		ID:        website.ID,
		Name:      website.Name,
		Domain:    website.Domain,
		CreatedAt: website.CreatedAt,
	}
}

func WebsitesFromDomain(websites []domain.Website) []Website {
	result := make([]Website, 0, len(websites))
	for _, website := range websites {
		result = append(result, WebsiteFromDomain(website))
	}

	return result
}

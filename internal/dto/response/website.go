package response

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

type WebsiteListOutput struct {
	Body []Website
}

type WebsiteOutput struct {
	Body Website
}

func NewWebsiteListOutput(websites []Website) *WebsiteListOutput {
	return &WebsiteListOutput{Body: websites}
}

func NewWebsiteOutput(website Website) *WebsiteOutput {
	return &WebsiteOutput{Body: website}
}

func ToWebsite(website domain.Website) Website {
	return Website{
		ID:        website.ID,
		Name:      website.Name,
		Domain:    website.Domain,
		CreatedAt: website.CreatedAt,
	}
}

func ToWebsites(websites []domain.Website) []Website {
	result := make([]Website, 0, len(websites))
	for _, website := range websites {
		result = append(result, ToWebsite(website))
	}

	return result
}

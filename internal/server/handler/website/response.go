package website

import (
	"time"

	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/website"
)

type Website struct {
	ID        string    `json:"id" format:"uuid"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"createdAt"`
}

type listOutput struct {
	Body []Website
}

type websiteOutput struct {
	Body Website
}

type WebsiteOK struct {
	OK bool `json:"ok"`
}

type okOutput struct {
	Body WebsiteOK
}

func newListOutput(websites []Website) *listOutput {
	return &listOutput{Body: websites}
}

func newWebsiteOutput(website Website) *websiteOutput {
	return &websiteOutput{Body: website}
}

func newOKOutput() *okOutput {
	return &okOutput{Body: WebsiteOK{OK: true}}
}

func ToWebsite(website websitedomain.Website) Website {
	return Website{
		ID:        website.ID,
		Name:      website.Name,
		Domain:    website.Domain,
		CreatedAt: website.CreatedAt,
	}
}

func ToWebsites(websites []websitedomain.Website) []Website {
	result := make([]Website, 0, len(websites))
	for _, website := range websites {
		result = append(result, ToWebsite(website))
	}

	return result
}

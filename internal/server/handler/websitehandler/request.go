package websitehandler

type WebsiteRequest struct {
	Name   string `json:"name" required:"true" minLength:"1" maxLength:"100"`
	Domain string `json:"domain,omitempty" maxLength:"500"`
}

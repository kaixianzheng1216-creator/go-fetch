package service

import (
	"net/url"
	"strings"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/util"
)

type trackingURLs struct {
	page     *url.URL
	referrer *url.URL
}

type trackingURLFields struct {
	path           string
	query          string
	referrerPath   string
	referrerQuery  string
	referrerDomain string
	hostname       string
}

type utmFields struct {
	source   string
	medium   string
	campaign string
	content  string
	term     string
}

func newTrackingURLs(event domain.TrackedEvent, website domain.Website) trackingURLs {
	pageURL := parsePageURL(event.URL, websiteFallbackHost(website))
	return trackingURLs{
		page:     pageURL,
		referrer: parseReferrerURL(event.Referrer, pageURL),
	}
}

func newTrackingURLFields(trackingURLs trackingURLs) trackingURLFields {
	return trackingURLFields{
		path:           util.TruncateRunes(pathWithHash(trackingURLs.page), maxURLPartLength),
		query:          util.TruncateRunes(trackingURLs.page.RawQuery, maxURLPartLength),
		referrerPath:   util.TruncateRunes(trackingURLs.referrer.Path, maxURLPartLength),
		referrerQuery:  util.TruncateRunes(trackingURLs.referrer.RawQuery, maxURLPartLength),
		referrerDomain: util.TruncateRunes(trimWWW(trackingURLs.referrer.Hostname()), maxURLPartLength),
		hostname:       util.TruncateRunes(trackingURLs.page.Hostname(), maxHostnameLength),
	}
}

func parsePageURL(rawURL, fallbackHost string) *url.URL {
	base := url.URL{Scheme: defaultURLScheme, Host: fallbackHost}
	if rawURL == "" {
		return &url.URL{Scheme: base.Scheme, Host: base.Host, Path: defaultURLPath}
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return &url.URL{Scheme: base.Scheme, Host: base.Host, Path: defaultURLPath}
	}

	return base.ResolveReference(parsedURL)
}

func websiteFallbackHost(website domain.Website) string {
	domainName := strings.TrimSpace(website.Domain)
	if domainName == "" {
		return website.ID.String()
	}

	if !strings.Contains(domainName, "://") && !strings.HasPrefix(domainName, "//") {
		domainName = "//" + domainName
	}

	if parsedURL, err := url.Parse(domainName); err == nil && parsedURL.Host != "" {
		return parsedURL.Host
	}

	return website.ID.String()
}

func parseReferrerURL(rawURL string, pageURL *url.URL) *url.URL {
	if rawURL == "" {
		return &url.URL{}
	}

	base := url.URL{Scheme: defaultURLScheme, Host: pageURL.Host}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return &url.URL{}
	}

	return base.ResolveReference(parsedURL)
}

func pathWithHash(pageURL *url.URL) string {
	path := pageURL.EscapedPath()
	if path == "" {
		path = defaultURLPath
	}

	if pageURL.Fragment != "" {
		path += "#" + pageURL.EscapedFragment()
	}

	return path
}

func trimWWW(host string) string {
	return strings.TrimPrefix(host, "www.")
}

func newUTMFields(values url.Values) utmFields {
	return utmFields{
		source:   util.TruncateRunes(values.Get("utm_source"), maxUTMValueLength),
		medium:   util.TruncateRunes(values.Get("utm_medium"), maxUTMValueLength),
		campaign: util.TruncateRunes(values.Get("utm_campaign"), maxUTMValueLength),
		content:  util.TruncateRunes(values.Get("utm_content"), maxUTMValueLength),
		term:     util.TruncateRunes(values.Get("utm_term"), maxUTMValueLength),
	}
}

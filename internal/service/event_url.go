package service

import (
	"net/url"
	"strings"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/textutil"
)

type eventURLs struct {
	page     *url.URL
	referrer *url.URL
}

func newEventURLs(event domain.TrackedEvent, website domain.Website) eventURLs {
	pageURL := parsePageURL(event.URL, websiteFallbackHost(website))
	return eventURLs{
		page:     pageURL,
		referrer: parseReferrerURL(event.Referrer, pageURL),
	}
}

type eventURLFields struct {
	path           string
	query          string
	referrerPath   string
	referrerQuery  string
	referrerDomain string
	hostname       string
}

func newEventURLFields(eventURLs eventURLs) eventURLFields {
	return eventURLFields{
		path:           textutil.TruncateRunes(pathWithHash(eventURLs.page), maxURLPartLength),
		query:          textutil.TruncateRunes(eventURLs.page.RawQuery, maxURLPartLength),
		referrerPath:   textutil.TruncateRunes(eventURLs.referrer.Path, maxURLPartLength),
		referrerQuery:  textutil.TruncateRunes(eventURLs.referrer.RawQuery, maxURLPartLength),
		referrerDomain: textutil.TruncateRunes(trimWWW(eventURLs.referrer.Hostname()), maxURLPartLength),
		hostname:       textutil.TruncateRunes(eventURLs.page.Hostname(), maxHostnameLength),
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
	siteDomain := strings.TrimSpace(website.Domain)
	if siteDomain == "" {
		return website.ID.String()
	}

	if !strings.Contains(siteDomain, "://") && !strings.HasPrefix(siteDomain, "//") {
		siteDomain = "//" + siteDomain
	}

	if parsedURL, err := url.Parse(siteDomain); err == nil && parsedURL.Host != "" {
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

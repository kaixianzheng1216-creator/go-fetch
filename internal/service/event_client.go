package service

import (
	"strconv"
	"strings"

	"github.com/mileusna/useragent"
)

type eventClient struct {
	ip              string
	userAgent       string
	browser         string
	operatingSystem string
	device          string
}

func newEventClient(clientInfo ClientInfo, screen string) eventClient {
	browser, operatingSystem, device := parseUserAgent(clientInfo.UserAgent, screen)
	return eventClient{
		ip:              clientInfo.IP,
		userAgent:       clientInfo.UserAgent,
		browser:         browser,
		operatingSystem: operatingSystem,
		device:          device,
	}
}

func parseUserAgent(userAgentValue, screen string) (browser, osName, device string) {
	agent := useragent.Parse(userAgentValue)

	browser = agent.Name
	if browser == "" || agent.IsUnknown() {
		browser = "Unknown"
	}

	osName = agent.OS
	if osName == "" {
		osName = "Unknown"
	}

	switch {
	case agent.Mobile:
		device = "mobile"
	case agent.Tablet:
		device = "tablet"
	default:
		device = "desktop"
		if width, _, hasHeight := strings.Cut(screen, "x"); hasHeight {
			if screenWidth, err := strconv.Atoi(width); err == nil && screenWidth <= laptopMaxScreenWidth {
				device = "laptop"
			}
		}
	}

	return browser, osName, device
}

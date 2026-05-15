package useragent

import (
	"strconv"
	"strings"

	"github.com/mileusna/useragent"
)

const laptopMaxScreenWidth = 1280

func Parse(userAgentValue, screen string) (browser, osName, device string) {
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

func IsBot(userAgentValue string) bool {
	return useragent.Parse(userAgentValue).Bot
}

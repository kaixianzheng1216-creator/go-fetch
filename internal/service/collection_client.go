package service

import (
	"strconv"
	"strings"

	"github.com/mileusna/useragent"
)

type trackingClient struct {
	ip        string
	userAgent string
	browser   string
	os        string
	device    string
	bot       bool
	country   string
	region    string
	city      string
}

func newTrackingClient(clientInfo ClientInfo, screen string) trackingClient {
	agent := useragent.Parse(clientInfo.UserAgent)
	return trackingClient{
		ip:        clientInfo.IP,
		userAgent: clientInfo.UserAgent,
		browser:   browserName(agent),
		os:        operatingSystemName(agent),
		device:    deviceType(agent, screen),
		bot:       agent.Bot,
		country:   clientInfo.Country,
		region:    clientInfo.Region,
		city:      clientInfo.City,
	}
}

func browserName(agent useragent.UserAgent) string {
	browser := agent.Name
	if browser == "" || agent.IsUnknown() {
		return "Unknown"
	}
	return browser
}

func operatingSystemName(agent useragent.UserAgent) string {
	if agent.OS == "" {
		return "Unknown"
	}
	return agent.OS
}

func deviceType(agent useragent.UserAgent, screen string) string {
	switch {
	case agent.Mobile:
		return "mobile"
	case agent.Tablet:
		return "tablet"
	}

	if width, ok := screenWidth(screen); ok && width <= laptopMaxScreenWidth {
		return "laptop"
	}

	return "desktop"
}

func screenWidth(screen string) (int, bool) {
	width, _, hasHeight := strings.Cut(screen, "x")
	if !hasHeight {
		return 0, false
	}

	value, err := strconv.Atoi(width)
	if err != nil {
		return 0, false
	}

	return value, true
}

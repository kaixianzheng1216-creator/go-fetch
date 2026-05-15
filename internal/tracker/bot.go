package tracker

import "github.com/mileusna/useragent"

func IsBot(userAgent string) bool {
	return useragent.Parse(userAgent).Bot
}

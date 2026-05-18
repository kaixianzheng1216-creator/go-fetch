package service

import "time"

type clock func() time.Time

func systemClock() time.Time {
	return time.Now()
}

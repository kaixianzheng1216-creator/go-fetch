package domain

import (
	"slices"
	"time"
)

type DateUnit string

const (
	DateUnitHour  DateUnit = "hour"
	DateUnitDay   DateUnit = "day"
	DateUnitMonth DateUnit = "month"

	DefaultDateUnit     = DateUnitHour
	DefaultDateLookback = 24 * time.Hour
)

var dateUnits = [...]DateUnit{
	DateUnitHour,
	DateUnitDay,
	DateUnitMonth,
}

func ParseDateUnit(value string) (DateUnit, bool) {
	if value == "" {
		return DefaultDateUnit, true
	}

	dateUnit := DateUnit(value)
	if slices.Contains(dateUnits[:], dateUnit) {
		return dateUnit, true
	}

	return "", false
}

func DateUnitValues() []string {
	return stringEnumValues(dateUnits[:])
}

func DateTruncUnit(unit DateUnit) string {
	return string(NormalizeDateUnit(unit))
}

func FormatBucket(bucketTime time.Time, unit DateUnit) string {
	switch NormalizeDateUnit(unit) {
	case DateUnitMonth:
		return bucketTime.Format("2006-01")
	case DateUnitDay:
		return bucketTime.Format("01-02")
	default:
		return bucketTime.Format("15:04")
	}
}

func NormalizeDateUnit(unit DateUnit) DateUnit {
	parsedUnit, isSupportedDateUnit := ParseDateUnit(string(unit))
	if !isSupportedDateUnit {
		return DefaultDateUnit
	}

	return parsedUnit
}

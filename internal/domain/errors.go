package domain

import "errors"

var (
	// ErrNotFound indicates a requested domain entity does not exist.
	ErrNotFound = errors.New("not found")
	// ErrUnsupportedDateUnit indicates an unknown analytics date unit.
	ErrUnsupportedDateUnit = errors.New("unsupported date unit")
	// ErrUnsupportedMetricType indicates an unknown analytics metric type.
	ErrUnsupportedMetricType = errors.New("unsupported metric type")
)

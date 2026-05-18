package domain

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrUnsupportedDateUnit   = errors.New("unsupported date unit")
	ErrUnsupportedMetricType = errors.New("unsupported metric type")
)

package errorz

import "errors"

var (
	ErrUnsupportedMetricType = errors.New("unsupported metric type")
	ErrInvalidMetricValue    = errors.New("invalid metric value")
	ErrInvalidMetricName     = errors.New("invalid metric name")
	ErrUpdateMetricFailed    = errors.New("update metric failed")
	ErrSendMetric            = errors.New("send metric failed")
)

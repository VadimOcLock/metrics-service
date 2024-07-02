package errorz

import "errors"

var (
	ErrUnsupportedMetricType = errors.New("unsupported metric type")
	ErrInvalidMetricValue    = errors.New("invalid metric value")
	ErrInvalidMetricName     = errors.New("invalid metric name")
	ErrUpdateMetricFailed    = errors.New("update metric failed")
	ErrSendMetric            = errors.New("send metric failed")
	ErrSendMetricStatusNotOK = errors.New("send metric status not ok")
)

const (
	ErrMsgOnlyPOSTMethodAccept = "only POST method accept"
	ErrMsgEmptyMetricParam     = "empty metric param"
)

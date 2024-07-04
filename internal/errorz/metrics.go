package errorz

import "errors"

var (
	ErrUndefinedMetricType    = errors.New("undefined metric type")
	ErrUndefinedMetricName    = errors.New("no metric with this name")
	ErrInvalidMetricValue     = errors.New("invalid metric value")
	ErrInvalidMetricName      = errors.New("invalid metric name")
	ErrUpdateMetricFailed     = errors.New("update metric failed")
	ErrSendMetricStatusNotOK  = errors.New("send metric status not ok")
	ErrCantConvertAnyToString = errors.New("cant convert any to string")
)

const (
	ErrMsgOnlyPOSTMethodAccept = "only POST method accept"
	ErrMsgEmptyMetricParam     = "empty metric param"
	ErrMsgFindAllMetrics       = "find all metrics error"
	ErrMsgFindMetric           = "find metric error"
)

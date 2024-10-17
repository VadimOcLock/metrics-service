package errorz

import "errors"

var (
	ErrUndefinedMetricType     = errors.New("undefined metric type")
	ErrUndefinedMetricName     = errors.New("no metric with this name")
	ErrInvalidMetricValue      = errors.New("invalid metric value")
	ErrInvalidMetricName       = errors.New("invalid metric name")
	ErrUpdateMetricFailed      = errors.New("update metric failed")
	ErrSendMetricStatusNotOK   = errors.New("send metric status not ok")
	ErrCantConvertAnyToString  = errors.New("cant convert any to string")
	ErrInvalidAddressFormat    = errors.New("invalid address format, expected host:port")
	ErrGaugeTypeNilValue       = errors.New("value is nil for gauge type")
	ErrCounterTypeNilDelta     = errors.New("delta is nil for counter type")
	ErrIncorrectDatabaseSchema = errors.New("incorrect database schema")
	ErrNoSpecifiedDatabaseName = errors.New("no database name specified")
)

const (
	ErrMsgOnlyPOSTMethodAccept = "only POST method accept"
	ErrMsgEmptyMetricParam     = "empty metric param"
	ErrMsgFindAllMetrics       = "find all metrics error"
	ErrMsgFindMetric           = "find metric error"
)

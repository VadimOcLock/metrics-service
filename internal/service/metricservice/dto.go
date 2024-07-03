package metricservice

import (
	"github.com/VadimOcLock/metrics-service/internal/store/somestore"
)

type UpdateGaugeDTO somestore.UpdateGaugeMetricParams

func (dto UpdateGaugeDTO) Valid() error {
	// may be later

	return nil
}

type UpdateCounterDTO somestore.UpdateCounterMetricParams

func (dto UpdateCounterDTO) Valid() error {
	// may be later

	return nil
}

type FindAllDTO somestore.FindAllMetricsParams

func (dto FindAllDTO) Valid() error {
	// may be later

	return nil
}

type FindDTO struct {
	MetricType string
	MetricName string
}

func (dto FindDTO) Valid() error {
	// may be later

	return nil
}

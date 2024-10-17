package metricservice

import (
	"context"

	"github.com/VadimOcLock/metrics-service/internal/entity"
)

type Store interface {
	UpsertGaugeMetric(ctx context.Context, arg UpsertGaugeMetricParams) (bool, error)
	UpsertCounterMetric(ctx context.Context, arg UpsertCounterMetricParams) (bool, error)
	UpdateMetricsBatchTx(ctx context.Context, arg UpdateMetricsBatchTxParams) error
	FindGaugeMetrics(ctx context.Context, arg FindGaugeMetricParams) (entity.Metrics, error)
	FindCounterMetrics(ctx context.Context, arg FindCounterMetricParams) (entity.Metrics, error)
	FindAllMetrics(ctx context.Context, arg FindAllMetricsNewParams) ([]entity.Metrics, error)
}

type UpsertGaugeMetricParams struct {
	Name  string
	Value float64
}

type FindGaugeMetricParams struct {
	MetricName string
}

type FindCounterMetricParams struct {
	MetricName string
}

type UpsertCounterMetricParams struct {
	Name  string
	Value int64
}

type FindAllMetricsParams struct{}

type FindAllMetricsNewParams struct{}

type UpdateMetricsBatchTxParams struct {
	Data *[]entity.Metrics
}

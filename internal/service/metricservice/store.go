package metricservice

import (
	"context"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/store/somestore"
)

type Store interface {
	UpdateGaugeMetric(ctx context.Context, arg somestore.UpdateGaugeMetricParams) (bool, error)
	FindGaugeMetric(ctx context.Context, arg somestore.FindGaugeMetricParams) (entity.Metrics, error)
	FindCounterMetric(ctx context.Context, arg somestore.FindCounterMetricParams) (entity.Metrics, error)
	UpdateCounterMetric(ctx context.Context, arg somestore.UpdateCounterMetricParams) (bool, error)
	FindAllMetrics(ctx context.Context, arg somestore.FindAllMetricsParams) ([]entity.Metrics, error)
}

var _ Store = (*somestore.Impl)(nil)

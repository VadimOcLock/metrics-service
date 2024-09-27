package somestore

import (
	"context"
	"sync"

	"github.com/VadimOcLock/metrics-service/internal/entity"

	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
)

type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

type UpdateGaugeMetricParams struct {
	Name  string
	Value float64
}

func (i *Impl) UpdateGaugeMetric(ctx context.Context, arg UpdateGaugeMetricParams) (bool, error) {
	i.s.mu.Lock()
	defer i.s.mu.Unlock()
	i.s.gauges[arg.Name] = arg.Value

	return true, nil
}

type UpdateCounterMetricParams struct {
	Name  string
	Value int64
}

func (i *Impl) UpdateCounterMetric(ctx context.Context, arg UpdateCounterMetricParams) (bool, error) {
	i.s.mu.Lock()
	defer i.s.mu.Unlock()
	i.s.counters[arg.Name] += arg.Value

	return true, nil
}

type FindAllMetricsParams struct {
}

func (i *Impl) FindAllMetrics(_ context.Context, _ FindAllMetricsParams) ([]entity.Metric, error) {
	i.s.mu.RLock()
	defer i.s.mu.RUnlock()
	var metrics []entity.Metric
	for n, v := range i.s.gauges {
		metrics = append(metrics, entity.Metric{
			Type:  enum.GaugeMetricType,
			Name:  n,
			Value: v,
		})
	}
	for n, v := range i.s.counters {
		metrics = append(metrics, entity.Metric{
			Type:  enum.CounterMetricType,
			Name:  n,
			Value: v,
		})
	}

	return metrics, nil
}

type FindCounterMetricParams struct {
	MetricName string
}

func (i *Impl) FindCounterMetric(_ context.Context, arg FindCounterMetricParams) (entity.Metric, error) {
	i.s.mu.RLock()
	defer i.s.mu.RUnlock()
	metricValue, ok := i.s.counters[arg.MetricName]
	if !ok {
		return entity.Metric{}, errorz.ErrUndefinedMetricName
	}

	return entity.Metric{
		Type:  enum.CounterMetricType,
		Name:  arg.MetricName,
		Value: metricValue,
	}, nil
}

type FindGaugeMetricParams struct {
	MetricName string
}

func (i *Impl) FindGaugeMetric(_ context.Context, arg FindGaugeMetricParams) (entity.Metric, error) {
	i.s.mu.RLock()
	defer i.s.mu.RUnlock()
	metricValue, ok := i.s.gauges[arg.MetricName]
	if !ok {
		return entity.Metric{}, errorz.ErrUndefinedMetricName
	}

	return entity.Metric{
		Type:  enum.GaugeMetricType,
		Name:  arg.MetricName,
		Value: metricValue,
	}, nil
}

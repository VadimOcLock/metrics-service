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

func (i *Impl) FindAllMetrics(_ context.Context, _ FindAllMetricsParams) ([]entity.Metrics, error) {
	i.s.mu.Lock()
	defer i.s.mu.Unlock()
	var metrics []entity.Metrics
	for n, v := range i.s.gauges {
		metrics = append(metrics, entity.Metrics{
			MType: enum.GaugeMetricType,
			ID:    n,
			Value: &v,
		})
	}
	for n, v := range i.s.counters {
		metrics = append(metrics, entity.Metrics{
			MType: enum.CounterMetricType,
			ID:    n,
			Delta: &v,
		})
	}

	return metrics, nil
}

type FindCounterMetricParams struct {
	MetricName string
}

func (i *Impl) FindCounterMetric(_ context.Context, arg FindCounterMetricParams) (entity.Metrics, error) {
	i.s.mu.Lock()
	defer i.s.mu.Unlock()
	metricValue, ok := i.s.counters[arg.MetricName]
	if !ok {
		return entity.Metrics{}, errorz.ErrUndefinedMetricName
	}

	return entity.Metrics{
		MType: enum.CounterMetricType,
		ID:    arg.MetricName,
		Delta: &metricValue,
	}, nil
}

type FindGaugeMetricParams struct {
	MetricName string
}

func (i *Impl) FindGaugeMetric(_ context.Context, arg FindGaugeMetricParams) (entity.Metrics, error) {
	i.s.mu.Lock()
	defer i.s.mu.Unlock()
	metricValue, ok := i.s.gauges[arg.MetricName]
	if !ok {
		return entity.Metrics{}, errorz.ErrUndefinedMetricName
	}

	return entity.Metrics{
		MType: enum.GaugeMetricType,
		ID:    arg.MetricName,
		Value: &metricValue,
	}, nil
}

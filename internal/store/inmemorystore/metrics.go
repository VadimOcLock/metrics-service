package inmemorystore

import (
	"context"
	"sync"

	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"

	"github.com/VadimOcLock/metrics-service/internal/entity"
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

func (i *Impl) UpsertGaugeMetric(_ context.Context, arg metricservice.UpsertGaugeMetricParams) (bool, error) {
	i.s.mu.Lock()
	defer i.s.mu.Unlock()
	i.s.gauges[arg.Name] = arg.Value

	return true, nil
}

type UpdateCounterMetricParams struct {
	Name  string
	Value int64
}

func (i *Impl) UpsertCounterMetric(_ context.Context, arg metricservice.UpsertCounterMetricParams) (bool, error) {
	i.s.mu.Lock()
	defer i.s.mu.Unlock()
	i.s.counters[arg.Name] += arg.Value

	return true, nil
}

func (i *Impl) FindGaugeMetrics(ctx context.Context, arg metricservice.FindGaugeMetricParams) (entity.Metrics, error) {
	i.s.mu.RLock()
	defer i.s.mu.RUnlock()
	vl, ok := i.s.gauges[arg.MetricName]
	if !ok {
		return entity.Metrics{}, errorz.ErrUndefinedMetricName
	}

	return entity.Metrics{
		ID:    enum.GaugeMetricType,
		MType: arg.MetricName,
		Value: &vl,
	}, nil
}

func (i *Impl) FindCounterMetrics(_ context.Context, arg metricservice.FindCounterMetricParams) (entity.Metrics, error) {
	i.s.mu.RLock()
	defer i.s.mu.RUnlock()
	vl, ok := i.s.counters[arg.MetricName]
	if !ok {
		return entity.Metrics{}, errorz.ErrUndefinedMetricName
	}

	return entity.Metrics{
		ID:    enum.CounterMetricType,
		MType: arg.MetricName,
		Delta: &vl,
	}, nil
}

func (i *Impl) FindAllMetrics(_ context.Context, _ metricservice.FindAllMetricsNewParams) ([]entity.Metrics, error) {
	i.s.mu.RLock()
	defer i.s.mu.RUnlock()
	res := make([]entity.Metrics, 0, len(i.s.gauges)+len(i.s.counters))
	for name, vl := range i.s.counters {
		res = append(res, entity.Metrics{
			ID:    name,
			MType: enum.CounterMetricType,
			Delta: &vl,
		})
	}
	for name, vl := range i.s.gauges {
		res = append(res, entity.Metrics{
			ID:    name,
			MType: enum.GaugeMetricType,
			Value: &vl,
		})
	}

	return res, nil
}

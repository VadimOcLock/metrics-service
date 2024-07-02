package somestore

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"sync"
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

func (i *Impl) FindAllMetrics(_ context.Context, _ FindAllMetricsParams) ([]Metric, error) {
	i.s.mu.Lock()
	defer i.s.mu.Unlock()
	var metrics []Metric
	for n, v := range i.s.gauges {
		metrics = append(metrics, Metric{
			Type:  enum.GaugeMetricType,
			Name:  n,
			Value: v,
		})
	}
	for n, v := range i.s.counters {
		metrics = append(metrics, Metric{
			Type:  enum.CounterMetricType,
			Name:  n,
			Value: v,
		})
	}

	return metrics, nil
}

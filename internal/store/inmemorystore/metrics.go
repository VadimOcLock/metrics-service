package inmemorystore

import (
	"context"

	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"

	"github.com/VadimOcLock/metrics-service/internal/entity"
)

func (i *Impl) UpsertGaugeMetric(_ context.Context, arg metricservice.UpsertGaugeMetricParams) (bool, error) {
	i.S.Mu.Lock()
	defer i.S.Mu.Unlock()
	i.S.Gauges[arg.Name] = arg.Value

	return true, nil
}

func (i *Impl) UpsertCounterMetric(_ context.Context, arg metricservice.UpsertCounterMetricParams) (bool, error) {
	i.S.Mu.Lock()
	defer i.S.Mu.Unlock()
	i.S.Counters[arg.Name] += arg.Value

	return true, nil
}

func (i *Impl) FindGaugeMetrics(ctx context.Context, arg metricservice.FindGaugeMetricParams) (entity.Metrics, error) {
	i.S.Mu.RLock()
	defer i.S.Mu.RUnlock()
	vl, ok := i.S.Gauges[arg.MetricName]
	if !ok {
		return entity.Metrics{}, errorz.ErrUndefinedMetricName
	}

	return entity.Metrics{
		ID:    arg.MetricName,
		MType: enum.GaugeMetricType,
		Value: &vl,
	}, nil
}

func (i *Impl) FindCounterMetrics(
	_ context.Context,
	arg metricservice.FindCounterMetricParams,
) (entity.Metrics, error) {
	i.S.Mu.RLock()
	defer i.S.Mu.RUnlock()
	vl, ok := i.S.Counters[arg.MetricName]
	if !ok {
		return entity.Metrics{}, errorz.ErrUndefinedMetricName
	}

	return entity.Metrics{
		ID:    arg.MetricName,
		MType: enum.CounterMetricType,
		Delta: &vl,
	}, nil
}

func (i *Impl) FindAllMetrics(_ context.Context, _ metricservice.FindAllMetricsNewParams) ([]entity.Metrics, error) {
	i.S.Mu.RLock()
	defer i.S.Mu.RUnlock()
	res := make([]entity.Metrics, 0, len(i.S.Gauges)+len(i.S.Counters))

	for name, vl := range i.S.Counters {
		vlCopy := vl
		res = append(res, entity.Metrics{
			ID:    name,
			MType: enum.CounterMetricType,
			Delta: &vlCopy,
		})
	}
	for name, vl := range i.S.Gauges {
		vlCopy := vl
		res = append(res, entity.Metrics{
			ID:    name,
			MType: enum.GaugeMetricType,
			Value: &vlCopy,
		})
	}

	return res, nil
}

func (i *Impl) UpdateMetricsBatchTx(ctx context.Context, arg metricservice.UpdateMetricsBatchTxParams) error {
	i.S.Mu.Lock()
	defer i.S.Mu.Unlock()
	metrics := *arg.Data
	for _, m := range metrics {
		switch m.MType {
		case enum.CounterMetricType:
			if m.Delta != nil {
				i.S.Counters[m.ID] += *m.Delta
			}
		case enum.GaugeMetricType:
			if m.Value != nil {
				i.S.Gauges[m.ID] = *m.Value
			}
		}
	}

	return nil
}

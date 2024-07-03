package metricservice

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"

	"github.com/VadimOcLock/metrics-service/internal/store/somestore"
)

func (s Service) UpdateGauge(ctx context.Context, dto UpdateGaugeDTO) error {
	if err := dto.Valid(); err != nil {
		return err
	}
	_, err := s.Store.UpdateGaugeMetric(ctx, somestore.UpdateGaugeMetricParams{
		Name:  dto.Name,
		Value: dto.Value,
	})

	return err
}

func (s Service) UpdateCounter(ctx context.Context, dto UpdateCounterDTO) error {
	if err := dto.Valid(); err != nil {
		return err
	}
	_, err := s.Store.UpdateCounterMetric(ctx, somestore.UpdateCounterMetricParams{
		Name:  dto.Name,
		Value: dto.Value,
	})

	return err
}

func (s Service) FindAll(ctx context.Context, dto FindAllDTO) ([]entity.Metric, error) {
	if err := dto.Valid(); err != nil {
		return nil, err
	}
	mm, err := s.Store.FindAllMetrics(ctx, somestore.FindAllMetricsParams{})
	if err != nil {
		return nil, err
	}
	res := make([]entity.Metric, len(mm))
	for i, m := range mm {
		res[i] = m.Entity()
	}

	return res, nil
}

func (s Service) Find(ctx context.Context, dto FindDTO) (entity.Metric, error) {
	if err := dto.Valid(); err != nil {
		return entity.Metric{}, err
	}
	var m entity.Metric
	switch dto.MetricType {
	case enum.GaugeMetricType:
		sm, err := s.Store.FindGaugeMetric(ctx, somestore.FindGaugeMetricParams{
			MetricName: dto.MetricName,
		})
		if err != nil {
			return entity.Metric{}, err
		}
		m = sm.Entity()
	case enum.CounterMetricType:
		sm, err := s.Store.FindCounterMetric(ctx, somestore.FindCounterMetricParams{
			MetricName: dto.MetricName,
		})
		if err != nil {
			return entity.Metric{}, err
		}
		m = sm.Entity()
	default:

		return entity.Metric{}, errorz.ErrUndefinedMetricType
	}

	return m, nil
}

package metricservice

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/entity"

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

package metricservice

import (
	"context"
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

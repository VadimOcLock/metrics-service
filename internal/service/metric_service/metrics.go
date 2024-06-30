package metric_service

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/store/some_store"
)

func (s Service) UpdateGauge(ctx context.Context, dto UpdateGaugeDTO) error {
	if err := dto.Valid(); err != nil {
		return err
	}
	_, err := s.Store.UpdateGaugeMetric(ctx, some_store.UpdateGaugeMetricParams{
		Name:  dto.Name,
		Value: dto.Value,
	})

	return err
}

func (s Service) UpdateCounter(ctx context.Context, dto UpdateCounterDTO) error {
	if err := dto.Valid(); err != nil {
		return err
	}
	_, err := s.Store.UpdateCounterMetric(ctx, some_store.UpdateCounterMetricParams{
		Name:  dto.Name,
		Value: dto.Value,
	})

	return err
}

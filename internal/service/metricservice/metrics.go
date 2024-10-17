package metricservice

import (
	"context"
	"fmt"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
)

func (s *Service) UpdateGauge(ctx context.Context, dto UpdateGaugeDTO) error {
	if err := dto.Valid(); err != nil {
		return fmt.Errorf("metricservice.UpdateGauge: %w", err)
	}
	if _, err := s.Store.UpsertGaugeMetric(ctx, UpsertGaugeMetricParams{
		Name:  dto.Name,
		Value: dto.Value,
	}); err != nil {
		return fmt.Errorf("metricservice.UpdateGauge: %w", err)
	}

	return nil
}

func (s *Service) UpdateCounter(ctx context.Context, dto UpdateCounterDTO) error {
	if err := dto.Valid(); err != nil {
		return fmt.Errorf("metricservice.UpdateCounter: %w", err)
	}
	if _, err := s.Store.UpsertCounterMetric(ctx, UpsertCounterMetricParams{
		Name:  dto.Name,
		Value: dto.Value,
	}); err != nil {
		return fmt.Errorf("metricservice.UpdateCounter: %w", err)
	}

	return nil
}

func (s *Service) FindAll(ctx context.Context, dto FindAllDTO) ([]entity.Metrics, error) {
	if err := dto.Valid(); err != nil {
		return nil, fmt.Errorf("metricservice.FindAll: %w", err)
	}
	res, err := s.Store.FindAllMetrics(ctx, FindAllMetricsNewParams{})
	if err != nil {
		return nil, fmt.Errorf("metricservice.FindAll: %w", err)
	}

	return res, nil
}

func (s *Service) Find(ctx context.Context, dto FindDTO) (entity.Metrics, error) {
	if err := dto.Valid(); err != nil {
		return entity.Metrics{}, fmt.Errorf("metricservice.Find: %w", err)
	}
	var m entity.Metrics
	switch dto.MetricType {
	case enum.GaugeMetricType:
		sm, err := s.Store.FindGaugeMetrics(ctx, FindGaugeMetricParams{
			MetricName: dto.MetricName,
		})
		if err != nil {
			return entity.Metrics{}, err
		}
		m = sm
	case enum.CounterMetricType:
		sm, err := s.Store.FindCounterMetrics(ctx, FindCounterMetricParams{
			MetricName: dto.MetricName,
		})
		if err != nil {
			return entity.Metrics{}, err
		}
		m = sm
	default:

		return entity.Metrics{}, errorz.ErrUndefinedMetricType
	}
	if _, err := m.MetricValue(); err != nil {
		return entity.Metrics{}, err
	}

	return m, nil
}

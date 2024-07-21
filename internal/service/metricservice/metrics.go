package metricservice

import (
	"context"
	"fmt"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"

	"github.com/VadimOcLock/metrics-service/internal/store/somestore"
)

func (s *Service) UpdateGauge(ctx context.Context, dto UpdateGaugeDTO) error {
	if err := dto.Valid(); err != nil {
		return fmt.Errorf("metricservice.UpdateGauge: %w", err)
	}
	if _, err := s.Store.UpdateGaugeMetric(ctx, somestore.UpdateGaugeMetricParams{
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
	if _, err := s.Store.UpdateCounterMetric(ctx, somestore.UpdateCounterMetricParams{
		Name:  dto.Name,
		Value: dto.Value,
	}); err != nil {
		return fmt.Errorf("metricservice.UpdateCounter: %w", err)
	}

	return nil
}

func (s *Service) FindAll(ctx context.Context, dto FindAllDTO) ([]entity.Metric, error) {
	if err := dto.Valid(); err != nil {
		return nil, fmt.Errorf("metricservice.FindAll: %w", err)
	}
	res, err := s.Store.FindAllMetrics(ctx, somestore.FindAllMetricsParams{})
	if err != nil {
		return nil, fmt.Errorf("metricservice.FindAll: %w", err)
	}

	return res, nil
}

func (s *Service) Find(ctx context.Context, dto FindDTO) (entity.Metric, error) {
	if err := dto.Valid(); err != nil {
		return entity.Metric{}, fmt.Errorf("metricservice.Find: %w", err)
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
		m = sm
	case enum.CounterMetricType:
		sm, err := s.Store.FindCounterMetric(ctx, somestore.FindCounterMetricParams{
			MetricName: dto.MetricName,
		})
		if err != nil {
			return entity.Metric{}, err
		}
		m = sm
	default:

		return entity.Metric{}, errorz.ErrUndefinedMetricType
	}

	return m, nil
}

package metricusecase

import (
	"context"
	"fmt"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"

	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
)

func (uc *MetricUseCase) Update(ctx context.Context, dto MetricUpdateDTO) (entity.Metrics, error) {
	switch dto.MType {
	case enum.GaugeMetricType:
		if dto.Value == nil {
			return entity.Metrics{}, errorz.ErrInvalidMetricValue
		}
		err := uc.metricService.UpdateGauge(ctx, metricservice.UpdateGaugeDTO{
			Name:  dto.ID,
			Value: *dto.Value,
		})
		if err != nil {
			return entity.Metrics{}, errorz.ErrUpdateMetricFailed
		}
	case enum.CounterMetricType:
		if dto.Delta == nil {
			return entity.Metrics{}, errorz.ErrInvalidMetricValue
		}
		err := uc.metricService.UpdateCounter(ctx, metricservice.UpdateCounterDTO{
			Name:  dto.ID,
			Value: *dto.Delta,
		})
		if err != nil {
			return entity.Metrics{}, errorz.ErrUpdateMetricFailed
		}
	default:
		return entity.Metrics{}, errorz.ErrUndefinedMetricType
	}

	return entity.Metrics(dto), nil
}

func (uc *MetricUseCase) FindAll(ctx context.Context, _ MetricFindAllDTO) (MetricFindAllResp, error) {
	metrics, err := uc.metricService.FindAll(ctx, metricservice.FindAllDTO{})
	if err != nil {
		return MetricFindAllResp{}, fmt.Errorf("metricusecase.FindAll: %w", err)
	}
	html, err := buildHTML(metrics)
	if err != nil {
		return MetricFindAllResp{}, fmt.Errorf("metricusecase.FindAll: %w", err)
	}

	return MetricFindAllResp{
		HTML: html,
	}, nil
}

func (uc *MetricUseCase) Find(ctx context.Context, dto MetricFindDTO) (entity.Metrics, error) {
	metrics, err := uc.metricService.Find(ctx, metricservice.FindDTO{
		MetricType: dto.MetricType,
		MetricName: dto.MetricName,
	})
	if err != nil {
		return entity.Metrics{}, fmt.Errorf("metricusecase.Find: %w", err)
	}

	return metrics, nil
}

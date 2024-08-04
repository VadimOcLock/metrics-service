package metricusecase

import (
	"context"
	"fmt"
	"github.com/VadimOcLock/metrics-service/internal/entity"
	"strconv"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"

	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
)

func (uc *MetricUseCase) Update(ctx context.Context, dto MetricUpdateDTO) (MetricUpdateResp, error) {
	switch dto.Type {
	case enum.GaugeMetricType:
		vl, err := strconv.ParseFloat(dto.Value, 64)
		if err != nil {
			return MetricUpdateResp{}, errorz.ErrInvalidMetricValue
		}
		if dto.Name == "" {
			return MetricUpdateResp{}, errorz.ErrInvalidMetricName
		}
		err = uc.metricService.UpdateGauge(ctx, metricservice.UpdateGaugeDTO{
			Name:  dto.Name,
			Value: vl,
		})
		if err != nil {
			return MetricUpdateResp{}, errorz.ErrUpdateMetricFailed
		}
	case enum.CounterMetricType:
		vl, err := strconv.ParseInt(dto.Value, 10, 64)
		if err != nil {
			return MetricUpdateResp{}, errorz.ErrInvalidMetricValue
		}
		if dto.Name == "" {
			return MetricUpdateResp{}, errorz.ErrInvalidMetricName
		}
		err = uc.metricService.UpdateCounter(ctx, metricservice.UpdateCounterDTO{
			Name:  dto.Name,
			Value: vl,
		})
		if err != nil {
			return MetricUpdateResp{}, errorz.ErrUpdateMetricFailed
		}
	default:
		return MetricUpdateResp{}, errorz.ErrUndefinedMetricType
	}

	metrics, err := entity.BuildMetrics(entity.MetricDTO(dto))
	if err != nil {
		return MetricUpdateResp{}, fmt.Errorf("build err: %w", err)
	}

	return MetricUpdateResp{
		Message: "metric update success",
		Data:    &metrics,
	}, nil
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

func (uc *MetricUseCase) Find(ctx context.Context, dto MetricFindDTO) (MetricFindResp, error) {
	m, err := uc.metricService.Find(ctx, metricservice.FindDTO{
		MetricType: dto.MetricType,
		MetricName: dto.MetricName,
	})
	if err != nil {
		return MetricFindResp{}, fmt.Errorf("metricusecase.Find: %w", err)
	}

	metrics, err := entity.BuildMetrics(entity.MetricDTO{
		Type:  dto.MetricType,
		Name:  dto.MetricName,
		Value: fmt.Sprintf("%v", m.Value),
	})
	if err != nil {
		return MetricFindResp{}, fmt.Errorf("metricusecase.Find: %w", err)
	}

	return MetricFindResp{
		MetricValue: fmt.Sprintf("%v", m.Value),
		Data:        &metrics,
	}, nil
}

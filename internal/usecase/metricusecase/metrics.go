package metricusecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
)

func (uc MetricUseCase) Update(ctx context.Context, dto entity.MetricDTO) (UpdateResp, error) {

	switch dto.Type {
	case enum.GaugeMetricType:
		vl, err := strconv.ParseFloat(dto.Value, 64)
		if err != nil {
			return UpdateResp{}, errorz.ErrInvalidMetricValue
		}
		if dto.Name == "" {
			return UpdateResp{}, errorz.ErrInvalidMetricName
		}
		err = uc.metricService.UpdateGauge(ctx, metricservice.UpdateGaugeDTO{
			Name:  dto.Name,
			Value: vl,
		})
		if err != nil {
			return UpdateResp{}, errorz.ErrUpdateMetricFailed
		}
	case enum.CounterMetricType:
		vl, err := strconv.ParseInt(dto.Value, 10, 64)
		if err != nil {
			return UpdateResp{}, errorz.ErrInvalidMetricValue
		}
		if dto.Name == "" {
			return UpdateResp{}, errorz.ErrInvalidMetricName
		}
		err = uc.metricService.UpdateCounter(ctx, metricservice.UpdateCounterDTO{
			Name:  dto.Name,
			Value: vl,
		})
		if err != nil {
			return UpdateResp{}, errorz.ErrUpdateMetricFailed
		}
	default:
		return UpdateResp{}, errorz.ErrUndefinedMetricType
	}

	return UpdateResp{
		Message: "metric update success",
	}, nil
}

func (uc MetricUseCase) FindAll(ctx context.Context, _ FindAllDTO) (FindAllResp, error) {
	metrics, err := uc.metricService.FindAll(ctx, metricservice.FindAllDTO{})
	if err != nil {
		return FindAllResp{}, err
	}
	html, err := buildHTML(metrics)
	if err != nil {
		return FindAllResp{}, err
	}

	return FindAllResp{
		HTML: html,
	}, err
}

func (uc MetricUseCase) Find(ctx context.Context, dto FindDTO) (FindResp, error) {
	m, err := uc.metricService.Find(ctx, metricservice.FindDTO{
		MetricType: dto.MetricType,
		MetricName: dto.MetricName,
	})
	if err != nil {
		return FindResp{}, err
	}

	return FindResp{
		MetricValue: fmt.Sprintf("%v", m.Value),
	}, err
}

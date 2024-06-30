package metric_usecase

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"strconv"
)

func (uc UseCase) UpdateMetric(ctx context.Context, dto entity.MetricDTO) (UpdateMetricResp, error) {

	switch dto.Type {
	case enum.GaugeMetricType:
		vl, err := strconv.ParseFloat(dto.Value, 64)
		if err != nil {
			return UpdateMetricResp{}, errorz.ErrInvalidMetricValue
		}
		if dto.Name == "" {
			return UpdateMetricResp{}, errorz.ErrInvalidMetricName
		}
		err = uc.metricService.UpdateGauge(ctx, metricservice.UpdateGaugeDTO{
			Name:  dto.Name,
			Value: vl,
		})
		if err != nil {
			return UpdateMetricResp{}, errorz.ErrUpdateMetricFailed
		}
	case enum.CounterMetricType:
		vl, err := strconv.ParseInt(dto.Value, 10, 64)
		if err != nil {
			return UpdateMetricResp{}, errorz.ErrInvalidMetricValue
		}
		if dto.Name == "" {
			return UpdateMetricResp{}, errorz.ErrInvalidMetricName
		}
		err = uc.metricService.UpdateCounter(ctx, metricservice.UpdateCounterDTO{
			Name:  dto.Name,
			Value: vl,
		})
		if err != nil {
			return UpdateMetricResp{}, errorz.ErrUpdateMetricFailed
		}
	default:
		return UpdateMetricResp{}, errorz.ErrUnsupportedMetricType
	}

	return UpdateMetricResp{
		Message: "metric update success",
	}, nil
}

package metricusecase

import (
	"context"

	"github.com/VadimOcLock/metrics-service/internal/entity"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
)

type MetricUseCase struct {
	metricService metricservice.Service
}

type UseCase interface {
	UpdateMetric(ctx context.Context, dto entity.MetricDTO) (UpdateMetricResp, error)
}

type MetricService interface {
	UpdateGauge(ctx context.Context, dto metricservice.UpdateGaugeDTO) error
	UpdateCounter(ctx context.Context, dto metricservice.UpdateCounterDTO) error
}

func New(
	metricService metricservice.Service,
) UseCase {
	return MetricUseCase{
		metricService: metricService,
	}
}

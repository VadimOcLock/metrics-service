package metric_usecase

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/service/metric_service"
)

type UseCase struct {
	metricService metric_service.Service
}

type MetricService interface {
	UpdateGauge(ctx context.Context, dto metric_service.UpdateGaugeDTO) error
	UpdateCounter(ctx context.Context, dto metric_service.UpdateCounterDTO) error
}

func New(
	metricService metric_service.Service,
) UseCase {
	return UseCase{
		metricService: metricService,
	}
}

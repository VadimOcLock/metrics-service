package metricusecase

import (
	"context"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
)

type UseCase struct {
	metricService metricservice.Service
}

type MetricService interface {
	UpdateGauge(ctx context.Context, dto metricservice.UpdateGaugeDTO) error
	UpdateCounter(ctx context.Context, dto metricservice.UpdateCounterDTO) error
}

func New(
	metricService metricservice.Service,
) UseCase {
	return UseCase{
		metricService: metricService,
	}
}

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
	Update(ctx context.Context, dto entity.MetricDTO) (UpdateResp, error)
	FindAll(ctx context.Context, dto FindAllDTO) (FindAllResp, error)
	//Find(ctx context.Context, dto FindDTO) (FindResp, error)
}

type MetricService interface {
	UpdateGauge(ctx context.Context, dto metricservice.UpdateGaugeDTO) error
	UpdateCounter(ctx context.Context, dto metricservice.UpdateCounterDTO) error
	FindAll(ctx context.Context, dto metricservice.FindAllDTO) ([]entity.Metric, error)
}

func New(
	metricService metricservice.Service,
) UseCase {
	return MetricUseCase{
		metricService: metricService,
	}
}

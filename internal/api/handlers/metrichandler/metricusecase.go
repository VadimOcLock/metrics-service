package metrichandler

import (
	"context"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
)

type MetricUseCase interface {
	Update(ctx context.Context, dto metricusecase.MetricUpdateDTO) (entity.Metrics, error)
	FindAll(ctx context.Context, _ metricusecase.MetricFindAllDTO) (metricusecase.MetricFindAllResp, error)
	Find(ctx context.Context, dto metricusecase.MetricFindDTO) (entity.Metrics, error)
}

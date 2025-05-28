package metrichandler

import (
	"context"

	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
)

// MetricUseCase определяет контракт для use case работы с метриками.
type MetricUseCase interface {
	Update(ctx context.Context, dto metricusecase.MetricUpdateDTO) (metricusecase.MetricUpdateResp, error)
	UpdateBatch(ctx context.Context, dto metricusecase.MetricsUpdateBatchDTO) error
	FindAllWithHTML(ctx context.Context, _ metricusecase.MetricFindAllDTO) (metricusecase.MetricFindAllResp, error)
	Find(ctx context.Context, dto metricusecase.MetricFindDTO) (metricusecase.MetricFindResp, error)
}

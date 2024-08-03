package metricusecase

import (
	"context"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
)

type MetricService interface {
	UpdateGauge(ctx context.Context, dto metricservice.UpdateGaugeDTO) error
	UpdateCounter(ctx context.Context, dto metricservice.UpdateCounterDTO) error
	FindAll(ctx context.Context, dto metricservice.FindAllDTO) ([]entity.Metrics, error)
	Find(ctx context.Context, dto metricservice.FindDTO) (entity.Metrics, error)
}

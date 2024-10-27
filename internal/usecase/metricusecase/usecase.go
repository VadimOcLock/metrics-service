package metricusecase

import (
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
)

type MetricUseCase struct {
	metricService MetricService
}

var _ MetricService = (*metricservice.Service)(nil)

func New(
	metricService MetricService,
) MetricUseCase {
	return MetricUseCase{
		metricService: metricService,
	}
}

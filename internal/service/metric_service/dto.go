package metric_service

import (
	"github.com/VadimOcLock/metrics-service/internal/store/some_store"
)

type UpdateGaugeDTO some_store.UpdateGaugeMetricParams

func (dto UpdateGaugeDTO) Valid() error {
	// may be later

	return nil
}

type UpdateCounterDTO some_store.UpdateCounterMetricParams

func (dto UpdateCounterDTO) Valid() error {
	// may be later

	return nil
}

package some_store

import (
	"context"
)

type Store interface {
	UpdateGaugeMetric(ctx context.Context, arg UpdateGaugeMetricParams) (bool, error)
	UpdateCounterMetric(ctx context.Context, arg UpdateCounterMetricParams) (bool, error)
}

type Impl struct {
	s MemStorage
}

var _ Store = (*Impl)(nil)

func New() Impl {
	return Impl{
		s: MemStorage{
			gauges:   make(map[string]float64),
			counters: make(map[string]int64),
		},
	}
}

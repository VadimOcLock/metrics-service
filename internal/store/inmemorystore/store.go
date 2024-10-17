package inmemorystore

import (
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
)

type Impl struct {
	s MemStorage
}

func New() *Impl {
	return &Impl{
		s: MemStorage{
			gauges:   make(map[string]float64),
			counters: make(map[string]int64),
		},
	}
}

var _ metricservice.Store = (*Impl)(nil)

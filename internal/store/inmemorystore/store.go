package inmemorystore

import (
	"sync"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
)

type MemStorage struct {
	Mu       sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

type Impl struct {
	S *MemStorage
}

type Option func(*Impl)

func WithMemStorage(s *MemStorage) Option {
	return func(i *Impl) {
		i.S = s
	}
}

func New(options ...Option) *Impl {
	defaultStorage := &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	impl := &Impl{
		S: defaultStorage,
	}

	for _, opt := range options {
		opt(impl)
	}

	return impl
}

var _ metricservice.Store = (*Impl)(nil)

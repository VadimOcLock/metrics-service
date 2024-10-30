package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
)

type MetricsWorker struct {
	Opts      MetricsWorkerOpts
	MetricsCh chan entity.MetricsData
}

var _ lifecycle.WorkerRunner = (*MetricsWorker)(nil)

type MetricsWorkerOpts struct {
	ServerAddr         string
	PoolInterval       time.Duration
	ReportInterval     time.Duration
	SecretSignatureKey string
	RateLimit          int
}

func NewMetricsWorker(opts MetricsWorkerOpts) *MetricsWorker {
	return &MetricsWorker{
		Opts:      opts,
		MetricsCh: make(chan entity.MetricsData, opts.RateLimit),
	}
}

func (w *MetricsWorker) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	// Collect system metrics.
	wg.Add(1)
	go func() {
		defer wg.Done()
		w.collectSystemMetricsLoop(ctx, errCh)
	}()

	//Collect runtime metrics.
	wg.Add(1)
	go func() {
		defer wg.Done()
		w.collectRuntimeMetricsLoop(ctx, errCh)
	}()

	// Send metrics.
	for i := 0; i < w.Opts.RateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.sendMetricsLoop(ctx, errCh)
		}()
	}

	// Error handling.
	go func() {
		for err := range errCh {
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Error().Msg(err.Error())
			}
		}
	}()

	wg.Wait()
	close(errCh)
	log.Debug().Msg("worker goroutine finished")

	return nil
}

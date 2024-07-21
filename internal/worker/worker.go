package worker

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
	"time"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
)

type MetricsWorker struct {
	Opts MetricsWorkerOpts
}

var _ lifecycle.WorkerRunner = (*MetricsWorker)(nil)

type MetricsWorkerOpts struct {
	ServerAddr     string
	PoolInterval   time.Duration
	ReportInterval time.Duration
}

func NewMetricsWorker(opts MetricsWorkerOpts) *MetricsWorker {
	return &MetricsWorker{
		Opts: opts,
	}
}

func (w *MetricsWorker) Run(ctx context.Context) error {
	var metrics entity.Metrics
	var wg sync.WaitGroup
	chanErr := make(chan error, 1)
	pollTimer := time.NewTimer(w.Opts.PoolInterval)
	reportTimer := time.NewTimer(w.Opts.ReportInterval)
	defer func() {
		pollTimer.Stop()
		reportTimer.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()

			return ctx.Err()
		case err := <-chanErr:
			if err != nil {
				log.Error().Msg(err.Error())
			}
		case <-pollTimer.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := w.collectMetrics(ctx, &metrics)
				if err != nil {
					chanErr <- fmt.Errorf("worker.run: %w", err)
				}
				log.Debug().Msg("collect metric success")
				pollTimer.Reset(w.Opts.PoolInterval)
			}()
		case <-reportTimer.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := w.sendMetrics(ctx, &metrics)
				if err != nil {
					chanErr <- fmt.Errorf("worker.run: %w", err)
				}
				log.Debug().Msg("send metric success")
				reportTimer.Reset(w.Opts.ReportInterval)
			}()
		}
	}
}

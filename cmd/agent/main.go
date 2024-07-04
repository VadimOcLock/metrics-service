package main

import (
	"context"
	"os"
	"syscall"

	"github.com/VadimOcLock/metrics-service/internal/worker"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
)

func main() {
	ctx := context.Background()
	parseFlags()

	w := worker.NewMetricsWorker(worker.MetricsWorkerOpts{
		ServerAddr:     flagOpts.EndpointAddr.String(),
		PoolInterval:   flagOpts.PoolInterval,
		ReportInterval: flagOpts.ReportInterval,
	})

	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.Worker(w))
	_ = tasks.Run()
}

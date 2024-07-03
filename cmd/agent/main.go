package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/VadimOcLock/metrics-service/internal/worker"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
)

func main() {
	ctx := context.Background()
	parseFlags()

	w := worker.NewMetricsWorker(worker.MetricsWorkerOpts{
		ServerAddr:     flagRunAddr,
		PoolInterval:   time.Duration(flagPoolInterval) * time.Second,
		ReportInterval: time.Duration(flagReportInterval) * time.Second,
	})

	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.Worker(w))
	_ = tasks.Run()
}

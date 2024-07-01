package main

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/worker"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
	"os"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	srvAddr := "http://localhost:8080"

	w := worker.NewMetricsWorker(worker.MetricsWorkerOpts{
		ServerAddr:     srvAddr,
		PoolInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	})

	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.Worker(w))
	_ = tasks.Run()
}

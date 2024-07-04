package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/VadimOcLock/metrics-service/internal/config"

	"github.com/VadimOcLock/metrics-service/internal/worker"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load[config.Agent]()
	if err != nil {
		log.Println("cfg load err: ", err)
		os.Exit(1)
	}
	if err = parseFlags(&cfg); err != nil {
		log.Println("parse flags err: ", err)
		os.Exit(1)
	}

	w := worker.NewMetricsWorker(worker.MetricsWorkerOpts{
		ServerAddr:     HTTPProtocolName + "://" + cfg.EndpointAddr,
		PoolInterval:   time.Duration(cfg.AgentConfig.PoolInterval) * time.Second,
		ReportInterval: time.Duration(cfg.AgentConfig.ReportInterval) * time.Second,
	})

	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.Worker(w))
	_ = tasks.Run()
}

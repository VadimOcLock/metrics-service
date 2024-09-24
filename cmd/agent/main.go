package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/VadimOcLock/metrics-service/internal/config"

	"github.com/VadimOcLock/metrics-service/internal/worker"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
)

func main() {
	ctx := context.Background()

	// Config.
	cfg, err := config.Load[config.Agent]()
	if err != nil {
		log.Fatal().Msgf("cfg load err: %v", err)
	}

	// Flags.
	if err = parseFlags(&cfg); err != nil {
		log.Fatal().Msgf("parse flags err: %v", err)
	}

	// Logger.
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().Timestamp().Logger()

	// Worker.
	w := worker.NewMetricsWorker(worker.MetricsWorkerOpts{
		ServerAddr:     HTTPProtocolName + "://" + cfg.EndpointAddr,
		PoolInterval:   time.Duration(cfg.AgentConfig.PoolInterval) * time.Second,
		ReportInterval: time.Duration(cfg.AgentConfig.ReportInterval) * time.Second,
	})

	// Run app.
	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.Worker(w))
	if err = tasks.Run(); err != nil {
		log.Debug().Msgf("tasks shutdown err: %v", err)
	}
}

package main

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/store/somestore"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
	"github.com/VadimOcLock/metrics-service/internal/worker"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/metrichandler"

	"github.com/VadimOcLock/metrics-service/internal/config"

	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
)

func main() {
	ctx := context.Background()

	// Config.
	cfg, err := config.Load[config.WebServer]()
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

	// Store.
	store := somestore.New()

	// Service.
	metricService := metricservice.New(&store)

	// UseCase.
	metricUseCase := metricusecase.New(&metricService)

	// Handler.
	mh := metrichandler.NewMetricHandler(&metricUseCase)
	mux := metrichandler.New(mh)
	server := &http.Server{
		Addr:              cfg.WebServerConfig.SrvAddr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	// Backup worker.
	bw, err := worker.NewBackupWorker(&metricService, &metricUseCase, worker.MetricsBackupOpts{
		Restore:  cfg.BackupConfig.Restore,
		Interval: cfg.BackupConfig.Interval,
		Filepath: cfg.BackupConfig.FileStoragePath,
	})

	// Run app.
	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.Worker(bw))
	tasks.Add(lifecycle.HTTPServer(server))
	if err = tasks.Run(); err != nil {
		log.Debug().Msgf("tasks shutdown err: %v", err)
	}
}

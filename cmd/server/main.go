package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/VadimOcLock/metrics-service/internal/store/migrations"

	"github.com/VadimOcLock/metrics-service/internal/store/pgstore"
	"github.com/VadimOcLock/metrics-service/pkg/pg"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/store/inmemorystore"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
	"github.com/VadimOcLock/metrics-service/internal/worker"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/metrichandler"

	"github.com/VadimOcLock/metrics-service/internal/config"

	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrationsPath = "file://internal/store/migrations"
)

func main() {
	ctx := context.Background()

	// HandlerConfig.
	cfg, err := env.ParseAs[config.WebServer]()
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

	// Database pool.
	dbPool, err := pg.New(ctx, pg.Config{
		DSN: cfg.DatabaseConfig.DSN,
	})
	if err != nil && !cfg.DatabaseConfig.InMemoryMode() {
		log.Fatal().Msgf("database connect err: %v", err)
	}
	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()

	// Store.
	store, err := setupStore(cfg.InMemoryMode(), cfg.DSN, dbPool)
	if err != nil {
		log.Error().Msgf("init store failed: %v", err)

		return
	}

	// Service.
	metricService := metricservice.New(store)

	// UseCase.
	metricUseCase := metricusecase.New(&metricService)

	// Handler.
	mh := metrichandler.NewMetricHandler(&metricUseCase)
	mux := metrichandler.New(mh,
		metrichandler.WithDBPool(dbPool),
		metrichandler.WithSecretSignatureKey(cfg.SecretSignatureKey))
	server := &http.Server{
		Addr:              cfg.WebServerConfig.SrvAddr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	// Backup worker.
	bw, err := worker.NewBackupWorker(&metricUseCase, worker.MetricsBackupOpts{
		Restore:  cfg.BackupConfig.Restore,
		Interval: cfg.BackupConfig.Interval,
		Filepath: cfg.BackupConfig.FileStoragePath,
	})
	if err != nil {
		log.Error().Msgf("new backup worker err: %v", err)

		return
	}

	// Run app.
	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	if cfg.DatabaseConfig.InMemoryMode() {
		tasks.Add(lifecycle.Worker(bw))
	}
	tasks.Add(lifecycle.HTTPServer(server))
	if err = tasks.Run(); err != nil {
		log.Debug().Msgf("tasks shutdown err: %v", err)
	}
}

func setupStore(isMemoryMode bool, dsn string, dbPool *pgxpool.Pool) (metricservice.Store, error) {
	var store metricservice.Store
	if isMemoryMode {
		store = inmemorystore.New()
	} else {
		if err := migrations.Run(dsn, migrationsPath); err != nil {
			log.Error().Msgf("migrations err: %v", err)

			return nil, fmt.Errorf("migrations failed: %w", err)
		}
		store = pgstore.NewPgStore(dbPool)
	}

	return store, nil
}

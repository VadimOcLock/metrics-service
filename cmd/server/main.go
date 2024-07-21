package main

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

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
	log.Logger = zerolog.New(os.Stdout)

	// Handler.
	mux := metrichandler.New()
	server := &http.Server{
		Addr:              cfg.WebServerConfig.SrvAddr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	// Run app.
	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.HTTPServer(server))
	if err = tasks.Run(); err != nil {
		log.Fatal().Msgf("tasks shutdown err: %v", err)
	}
}

package main

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/VadimOcLock/metrics-service/internal/file"
	"github.com/VadimOcLock/metrics-service/internal/store/somestore"

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
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Store.
	store := somestore.New()

	// File writer listener channel.
	fileUpdateCh := make(chan bool, 2)
	listener := file.NewListener(fileUpdateCh, &cfg.FileWriter, &store)

	// Handler.
	mux := metrichandler.New(&store, fileUpdateCh)
	server := &http.Server{
		Addr:              cfg.WebServerConfig.SrvAddr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	// Run app.
	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.HTTPServer(server))
	tasks.Add(lifecycle.Listener(listener))
	if err = tasks.Run(); err != nil {
		log.Debug().Msgf("tasks shutdown err: %v", err)
	}
}

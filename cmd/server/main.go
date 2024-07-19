package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/metrichandler"

	"github.com/VadimOcLock/metrics-service/internal/config"

	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load[config.WebServer]()
	if err != nil {
		log.Fatalf("cfg load err: %v", err)
	}
	if err = parseFlags(&cfg); err != nil {
		log.Fatalf("parse flags err: %v", err)
	}

	mux := metrichandler.New()
	server := &http.Server{
		Addr:              cfg.WebServerConfig.SrvAddr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.HTTPServer(server))
	if err = tasks.Run(); err != nil {
		log.Printf("tasks shutdown err: %v", err)
	}
}

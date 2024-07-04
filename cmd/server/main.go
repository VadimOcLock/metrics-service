package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/VadimOcLock/metrics-service/internal/config"

	"github.com/VadimOcLock/metrics-service/internal/api/handler"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load[config.WebServer]()
	if err != nil {
		log.Println("cfg load err: ", err)
		os.Exit(1)
	}
	if err = parseFlags(&cfg); err != nil {
		log.Println("parse flags err: ", err)
		os.Exit(1)
	}

	mux := handler.New()
	server := &http.Server{
		Addr:    cfg.WebServerConfig.SrvAddr,
		Handler: mux,
	}

	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.HTTPServer(server))
	_ = tasks.Run()
}

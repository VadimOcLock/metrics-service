package main

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/api/handler"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
	"net/http"
	"os"
	"syscall"
)

func main() {
	ctx := context.Background()
	addr := "localhost:8080"

	mux := handler.New(ctx)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.HTTPServer(server))
	_ = tasks.Run()
}

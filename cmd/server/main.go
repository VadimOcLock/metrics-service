package main

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/VadimOcLock/metrics-service/internal/api/handler"
	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"
	"github.com/safeblock-dev/wr/taskgroup"
)

func main() {
	ctx := context.Background()
	parseFlags()

	mux := handler.New()
	server := &http.Server{
		Addr:    flagOpts.SrvAddr.String(),
		Handler: mux,
	}

	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.HTTPServer(server))
	_ = tasks.Run()
}

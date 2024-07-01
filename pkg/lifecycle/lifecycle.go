package lifecycle

import (
	"context"
	"errors"
	"github.com/safeblock-dev/werr"
	"github.com/safeblock-dev/wr/taskgroup"
	"log"
	"net/http"
	"time"
)

const httpServerShutdownTTL = 10 * time.Second

func HTTPServer(server *http.Server) (taskgroup.ExecuteFn, taskgroup.InterruptFn) {
	execute := func() error {
		log.Println("HTTP server starting...")
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		log.Println("HTTP server finished")

		return err
	}
	interrupt := func(_ error) {
		ctx, cancel := context.WithTimeout(context.Background(), httpServerShutdownTTL)
		defer cancel()
		err := server.Shutdown(ctx)
		log.Println("HTTP server shutdown complete, err: ", err)
	}

	return execute, interrupt
}

type WorkerRunner interface {
	Run(ctx context.Context) error
}

func Worker(
	worker WorkerRunner,
) (taskgroup.ExecuteFn, taskgroup.InterruptFn) {
	ctx, cancel := context.WithCancel(context.Background())

	execute := func() error {
		log.Println("worker starting...")
		err := worker.Run(ctx)
		if errors.Is(err, context.Canceled) {
			err = nil
		}
		log.Println("worker finished")

		return werr.Wrap(err)
	}

	interrupt := func(_ error) {
		cancel()
		log.Println("worker shutdown complete")
	}

	return execute, interrupt
}

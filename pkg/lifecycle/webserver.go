package lifecycle

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/safeblock-dev/wr/taskgroup"
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

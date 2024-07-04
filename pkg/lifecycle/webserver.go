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
		log.Printf("HTTP server starting at addr: %s...", server.Addr)
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
		if err := server.Shutdown(ctx); err != nil {
			log.Println("shutdown server err: ", err)
		}
		log.Println("HTTP server shutdown complete")
	}

	return execute, interrupt
}

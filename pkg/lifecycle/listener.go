package lifecycle

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/safeblock-dev/werr"
	"github.com/safeblock-dev/wr/taskgroup"
)

type ListenerRunner interface {
	Run(ctx context.Context) error
}

func Listener(
	listener ListenerRunner,
) (taskgroup.ExecuteFn, taskgroup.InterruptFn) {
	ctx, cancel := context.WithCancel(context.Background())

	execute := func() error {
		log.Debug().Msg("listener starting...")
		err := listener.Run(ctx)
		if errors.Is(err, context.Canceled) {
			err = nil
		}
		log.Debug().Msg("listener finished")

		return werr.Wrap(err)
	}

	interrupt := func(_ error) {
		cancel()

		log.Debug().Msg("listener shutdown complete")
	}

	return execute, interrupt
}

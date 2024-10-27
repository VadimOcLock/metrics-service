package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	maxRetryAttempts = 3
	multiplier       = 2
)

func RunCtx(ctx context.Context, startDelay time.Duration, fn func(_ context.Context) error) error {
	var err error
	delay := startDelay
	for i := 0; i < maxRetryAttempts; i++ {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, delay)
		defer cancel()

		if err = fn(ctxWithTimeout); err == nil {
			return nil
		}

		log.Error().Msgf("Attempt %d failed: %v. Retrying in %s...", i+1, err, delay)

		delay *= multiplier

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("after %d attempts, last error: %w", maxRetryAttempts, err)
}

func Run(startDelay time.Duration, fn func() error) error {
	var err error
	delay := startDelay
	for i := 0; i < maxRetryAttempts; i++ {
		if err = fn(); err == nil {
			return nil
		}

		log.Error().Msgf("Attempt %d failed: %v. Retrying in %s...", i+1, err, delay)

		delay *= multiplier

		<-time.After(delay)
	}

	return fmt.Errorf("after %d attempts, last error: %w", maxRetryAttempts, err)
}

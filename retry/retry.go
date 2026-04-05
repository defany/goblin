package retry

import (
	"context"
	"math/rand/v2"
	"time"
)

func Do(ctx context.Context, fn func(context.Context) error, opts ...Option) error {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	var err error
	for attempt := range cfg.attempts {
		if err = fn(ctx); err == nil {
			return nil
		}

		if cfg.retryIf != nil && !cfg.retryIf(err) {
			return err
		}

		if attempt == cfg.attempts-1 {
			break
		}

		delay := backoff(cfg.baseDelay, cfg.maxDelay, attempt, cfg.jitter)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return err
}

func backoff(base, maxDelay time.Duration, attempt int, jitter bool) time.Duration {
	delay := base << attempt
	delay = min(delay, maxDelay)

	if jitter {
		delay = delay/2 + rand.N(delay/2)
	}

	return delay
}

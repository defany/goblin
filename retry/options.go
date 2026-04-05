package retry

import "time"

type config struct {
	attempts    int
	baseDelay   time.Duration
	maxDelay    time.Duration
	jitter      bool
	retryIf     func(error) bool
}

var defaultConfig = config{
	attempts:  3,
	baseDelay: 100 * time.Millisecond,
	maxDelay:  5 * time.Second,
	jitter:    true,
}

type Option func(*config)

func WithAttempts(n int) Option {
	return func(c *config) {
		c.attempts = n
	}
}

func WithBaseDelay(d time.Duration) Option {
	return func(c *config) {
		c.baseDelay = d
	}
}

func WithMaxDelay(d time.Duration) Option {
	return func(c *config) {
		c.maxDelay = d
	}
}

func WithoutJitter() Option {
	return func(c *config) {
		c.jitter = false
	}
}

func WithRetryIf(fn func(error) bool) Option {
	return func(c *config) {
		c.retryIf = fn
	}
}

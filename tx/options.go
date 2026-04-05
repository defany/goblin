package tx

import (
	"context"
	"time"
)

type Option func(*Manager)

func WithPanicHandler(fn func(ctx context.Context, p HandledPanic)) Option {
	return func(m *Manager) {
		m.panicHandler = fn
	}
}

func WithIsRetryable(fn func(error) bool) Option {
	return func(m *Manager) {
		m.isRetryable = fn
	}
}

type RunOption func(*runConfig)

type runConfig struct {
	iso        IsoLevel
	readOnly   bool
	retry      uint
	maxBackoff time.Duration
}

var defaultRunConfig = runConfig{
	retry: 5,
}

func WithIso(lvl IsoLevel) RunOption {
	return func(c *runConfig) {
		c.iso = lvl
	}
}

func WithReadOnly(on bool) RunOption {
	return func(c *runConfig) {
		c.readOnly = on
	}
}

func WithRetry(n uint) RunOption {
	return func(c *runConfig) {
		c.retry = n
	}
}

func WithMaxBackoff(d time.Duration) RunOption {
	return func(c *runConfig) {
		c.maxBackoff = d
	}
}

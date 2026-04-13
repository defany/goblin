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

func WithJoinRetryErrors(on bool) Option {
	return func(m *Manager) {
		m.joinRetryErrors = on
	}
}

type RunOption func(*runConfig)

type runConfig struct {
	iso         IsoLevel
	readOnly    bool
	retry       uint
	maxBackoff  time.Duration
	isRetryable func(error) bool
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

func WithIsErrorRetryable(fn func(error) bool) RunOption {
	return func(c *runConfig) {
		c.isRetryable = fn
	}
}

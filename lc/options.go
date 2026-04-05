package lc

import (
	"context"
	"log/slog"
	"time"
)

type Option func(*Lifecycle)

func WithContext(ctx context.Context) Option {
	return func(l *Lifecycle) {
		l.ctx, l.cancel = context.WithCancel(ctx)
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(l *Lifecycle) {
		l.logger = logger
	}
}

func WithShutdownTimeout(d time.Duration) Option {
	return func(l *Lifecycle) {
		l.shutdownTimeout = d
	}
}

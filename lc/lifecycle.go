package lc

import (
	"context"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sourcegraph/conc/pool"
)

type Lifecycle struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger *slog.Logger

	shutdownTimeout time.Duration

	mu     sync.Mutex
	defers []func(context.Context) error
	heads  []*Head
}

const defaultShutdownTimeout = 5 * time.Second

func New(opts ...Option) *Lifecycle {
	l := &Lifecycle{
		logger:          slog.Default(),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(l)
	}

	if l.ctx == nil {
		l.ctx, l.cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	} else {
		l.ctx, l.cancel = signal.NotifyContext(l.ctx, syscall.SIGINT, syscall.SIGTERM)
	}

	return l
}

func (l *Lifecycle) Context() context.Context {
	return with(l.ctx, l)
}

func (l *Lifecycle) Go(fn func(context.Context) error) *Head {
	h := &Head{l: l, fn: fn}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.heads = append(l.heads, h)
	return h
}

func (l *Lifecycle) addDefer(fn func(context.Context) error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.defers = append(l.defers, fn)
}

func (l *Lifecycle) Run() error {
	defer l.cancel()

	p := pool.New().WithErrors()

	l.mu.Lock()
	heads := make([]*Head, len(l.heads))
	copy(heads, l.heads)
	l.mu.Unlock()

	for _, h := range heads {
		p.Go(func() error {
			return h.run(l.ctx)
		})
	}

	err := p.Wait()

	l.mu.Lock()
	defers := make([]func(context.Context) error, len(l.defers))
	copy(defers, l.defers)
	l.mu.Unlock()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), l.shutdownTimeout)
	defer shutdownCancel()

	for i := len(defers) - 1; i >= 0; i-- {
		if dErr := defers[i](shutdownCtx); dErr != nil {
			l.logger.Error("defer failed", slog.String("error", dErr.Error()))
			if err == nil {
				err = dErr
			}
		}
	}

	return err
}

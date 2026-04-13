package tx

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"runtime/debug"
	"time"

	"github.com/defany/goblin/errfmt"
)

const baseBackoff = 10 * time.Millisecond

type Handler = func(context.Context) error

type HandledPanic struct {
	Err        error
	Stacktrace string
}

type ctxKey struct{}

func extractTx(ctx context.Context) Transaction {
	t, _ := ctx.Value(ctxKey{}).(Transaction)
	return t
}

func injectTx(ctx context.Context, t Transaction) context.Context {
	return context.WithValue(ctx, ctxKey{}, t)
}

type Manager struct {
	db              Beginner
	panicHandler    func(ctx context.Context, p HandledPanic)
	joinRetryErrors bool
}

func New(db Beginner, opts ...Option) *Manager {
	m := &Manager{db: db}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Manager) ReadCommitted(ctx context.Context, h Handler, opts ...RunOption) error {
	return m.Run(ctx, h, append(opts, WithIso(ReadCommittedIso))...)
}

func (m *Manager) RepeatableRead(ctx context.Context, h Handler, opts ...RunOption) error {
	return m.Run(ctx, h, append(opts, WithIso(RepeatableReadIso))...)
}

func (m *Manager) Serializable(ctx context.Context, h Handler, opts ...RunOption) error {
	return m.Run(ctx, h, append(opts, WithIso(SerializableIso))...)
}

func (m *Manager) Run(ctx context.Context, h Handler, opts ...RunOption) error {
	cfg := defaultRunConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	if extractTx(ctx) != nil {
		return m.execTx(ctx, cfg, h)
	}

	var errs []error

	for attempt := range cfg.retry + 1 {
		err := m.execTx(ctx, cfg, h)
		if err == nil {
			return nil
		}

		retryable := m.db.IsRetryable(err) || (cfg.isRetryable != nil && cfg.isRetryable(err))
		if !retryable {
			return err
		}

		errs = append(errs, err)

		if attempt < cfg.retry {
			delay := backoff(baseBackoff, cfg.maxBackoff, attempt)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	cause := errs[len(errs)-1]
	if m.joinRetryErrors {
		cause = errors.Join(errs...)
	}

	return errfmt.WithSource(fmt.Errorf("tx: retries exceeded (%d attempts): %w", len(errs), cause))
}

func (m *Manager) execTx(ctx context.Context, cfg runConfig, h Handler) (err error) {
	if extractTx(ctx) != nil {
		return h(ctx)
	}

	t, err := m.db.BeginTx(ctx, Options{
		IsoLevel: cfg.iso,
		ReadOnly: cfg.readOnly,
	})
	if err != nil {
		return errfmt.WithSource(fmt.Errorf("tx: begin: %w", err))
	}

	ctx = injectTx(ctx, t)
	ctx = t.InjectCtx(ctx)

	defer func() {
		if r := recover(); r != nil {
			err = errfmt.WithSource(fmt.Errorf("tx: panic: %v", r))

			if m.panicHandler != nil {
				m.panicHandler(ctx, HandledPanic{
					Err:        err,
					Stacktrace: string(debug.Stack()),
				})
			}
		}

		if err != nil {
			if rbErr := t.Rollback(ctx); rbErr != nil {
				err = errors.Join(err, rbErr)
			}

			return
		}

		err = t.Commit(ctx)
	}()

	return h(ctx)
}

func backoff(base, maxDelay time.Duration, attempt uint) time.Duration {
	delay := base << attempt
	if maxDelay > 0 {
		delay = min(delay, maxDelay)
	}
	return delay/2 + rand.N(delay/2)
}

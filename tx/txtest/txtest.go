// Package txtest provides a no-op transaction Manager for tests, so callers can
// exercise code that runs inside tx.Manager without a real database.
package txtest

import (
	"context"

	"github.com/defany/goblin/tx"
)

// NewManager returns a tx.Manager that executes handlers synchronously without a
// real database: BeginTx/Commit/Rollback are no-ops and the handler runs once.
func NewManager() *tx.Manager {
	return tx.New(noopBeginner{})
}

type noopBeginner struct{}

func (noopBeginner) BeginTx(context.Context, tx.Options) (tx.Transaction, error) {
	return noopTx{}, nil
}

func (noopBeginner) IsRetryable(error) bool { return false }

type noopTx struct{}

func (noopTx) Commit(context.Context) error                  { return nil }
func (noopTx) Rollback(context.Context) error                { return nil }
func (noopTx) InjectCtx(ctx context.Context) context.Context { return ctx }

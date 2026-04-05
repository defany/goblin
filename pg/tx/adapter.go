package pgtx

import (
	"context"
	"errors"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/goblin/pg"
	"github.com/defany/goblin/tx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	codeSerialization = "40001"
	codeDeadlock      = "40P01"
)

var isoMap = map[tx.IsoLevel]pgx.TxIsoLevel{
	tx.ReadCommittedIso:  pgx.ReadCommitted,
	tx.RepeatableReadIso: pgx.RepeatableRead,
	tx.SerializableIso:   pgx.Serializable,
}

type Adapter struct {
	db *pg.Postgres
}

func NewAdapter(db *pg.Postgres) *Adapter {
	return &Adapter{db: db}
}

func (a *Adapter) BeginTx(ctx context.Context, opts tx.Options) (tx.Transaction, error) {
	pgOpts := pgx.TxOptions{
		IsoLevel:   isoMap[opts.IsoLevel],
		AccessMode: accessMode(opts.ReadOnly),
	}

	t, err := a.db.BeginTx(ctx, pgOpts)
	if err != nil {
		return nil, errfmt.WithSource(err)
	}

	return &pgTx{tx: t}, nil
}

func (a *Adapter) IsRetryable(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == codeSerialization || pgErr.Code == codeDeadlock
}

type pgTx struct {
	tx pgx.Tx
}

func (t *pgTx) Commit(ctx context.Context) error {
	return errfmt.WithSource(t.tx.Commit(ctx))
}

func (t *pgTx) Rollback(ctx context.Context) error {
	return errfmt.WithSource(t.tx.Rollback(ctx))
}

func (t *pgTx) InjectCtx(ctx context.Context) context.Context {
	return pg.InjectTx(ctx, t.tx)
}

func accessMode(readOnly bool) pgx.TxAccessMode {
	if readOnly {
		return pgx.ReadOnly
	}

	return pgx.ReadWrite
}

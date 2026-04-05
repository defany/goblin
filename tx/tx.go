package tx

import "context"

type IsoLevel string

const (
	ReadCommittedIso  IsoLevel = "read_committed"
	RepeatableReadIso IsoLevel = "repeatable_read"
	SerializableIso   IsoLevel = "serializable"
)

type Options struct {
	IsoLevel IsoLevel
	ReadOnly bool
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	InjectCtx(ctx context.Context) context.Context
}

type Beginner interface {
	BeginTx(ctx context.Context, opts Options) (Transaction, error)
	IsRetryable(err error) bool
}

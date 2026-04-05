package tx

import "context"

type GenericHandler[T any] func(context.Context) (T, error)

func ReadCommitted[T any](ctx context.Context, m *Manager, h GenericHandler[T], opts ...RunOption) (T, error) {
	return run(ctx, m, h, append(opts, WithIso(ReadCommittedIso))...)
}

func RepeatableRead[T any](ctx context.Context, m *Manager, h GenericHandler[T], opts ...RunOption) (T, error) {
	return run(ctx, m, h, append(opts, WithIso(RepeatableReadIso))...)
}

func Serializable[T any](ctx context.Context, m *Manager, h GenericHandler[T], opts ...RunOption) (T, error) {
	return run(ctx, m, h, append(opts, WithIso(SerializableIso))...)
}

func run[T any](ctx context.Context, m *Manager, h GenericHandler[T], opts ...RunOption) (out T, err error) {
	err = m.Run(ctx, func(txCtx context.Context) error {
		out, err = h(txCtx)
		return err
	}, opts...)

	return
}

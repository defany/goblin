package lc

import "context"

type ctxKey struct{}
type readyKey struct{}

func with(ctx context.Context, l *Lifecycle) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func from(ctx context.Context) *Lifecycle {
	l, _ := ctx.Value(ctxKey{}).(*Lifecycle)
	return l
}

func Ready(ctx context.Context) {
	fn, _ := ctx.Value(readyKey{}).(func())
	if fn != nil {
		fn()
	}
}

func Defer(ctx context.Context, fn func(context.Context) error) {
	l := from(ctx)
	if l == nil {
		return
	}

	l.addDefer(fn)
}

func withReady(ctx context.Context, fn func()) context.Context {
	return context.WithValue(ctx, readyKey{}, fn)
}

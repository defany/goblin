package lc

import (
	"context"
	"sync"

	"github.com/sourcegraph/conc/pool"
)

type Head struct {
	l    *Lifecycle
	fn   func(context.Context) error
	then *Head
}

func (h *Head) Go(fn func(context.Context) error) *Head {
	next := &Head{l: h.l, fn: fn}
	h.then = next
	return next
}

func (h *Head) run(ctx context.Context) error {
	p := pool.New().WithErrors()

	lcCtx := with(ctx, h.l)

	for cur := h; cur != nil; cur = cur.then {
		ready := make(chan struct{})
		once := sync.Once{}

		fn := cur.fn
		next := cur.then

		p.Go(func() error {
			readyCtx := withReady(lcCtx, func() {
				once.Do(func() { close(ready) })
			})

			err := fn(readyCtx)
			once.Do(func() { close(ready) })
			return err
		})

		if next != nil {
			<-ready
		}
	}

	return p.Wait()
}

package inject

import (
	"context"
	"sync"

	"github.com/defany/goblin/rt"
	"golang.org/x/sync/singleflight"
)

var (
	onceStore     sync.Map
	singleRequest singleflight.Group
)

type onceEntry struct {
	val any
}

func Once[T any](ctx context.Context, f func(context.Context) T) T {
	key := rt.CallerUniqueKey(1)

	if val, ok := onceStore.Load(key); ok {
		return val.(*onceEntry).val.(T)
	}

	result, _, _ := singleRequest.Do(key, func() (any, error) {
		val := f(ctx)
		onceStore.Store(key, &onceEntry{val: val})
		return val, nil
	})

	return result.(T)
}

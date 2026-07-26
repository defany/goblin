package lc

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunWaitsHooksBeforeDefers(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	l := New(WithContext(parent))
	ctx := l.Context()

	var hookDone, deferAfterHook atomic.Bool

	OnShutdown(ctx, func(context.Context) error {
		time.Sleep(100 * time.Millisecond)
		hookDone.Store(true)

		return nil
	})

	Defer(ctx, func(context.Context) error {
		deferAfterHook.Store(hookDone.Load())

		return nil
	})

	l.Go(func(ctx context.Context) error {
		<-ctx.Done()

		return nil
	})

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	require.NoError(t, l.Run())
	require.True(t, hookDone.Load())
	require.True(t, deferAfterHook.Load())
}

func TestHeadErrorShutsDownOthers(t *testing.T) {
	l := New(WithContext(context.Background()))
	ctx := l.Context()

	wantErr := errors.New("boom")

	var hookRan atomic.Bool

	OnShutdown(ctx, func(context.Context) error {
		hookRan.Store(true)

		return nil
	})

	l.Go(func(context.Context) error {
		return wantErr
	})

	l.Go(func(ctx context.Context) error {
		<-ctx.Done()

		return nil
	})

	done := make(chan error, 1)
	go func() {
		done <- l.Run()
	}()

	select {
	case err := <-done:
		require.ErrorIs(t, err, wantErr)
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after head error")
	}

	require.True(t, hookRan.Load())
}

func TestHooksRunWhenHeadsFinishWithoutSignal(t *testing.T) {
	l := New(WithContext(context.Background()))
	ctx := l.Context()

	var hookRan atomic.Bool

	OnShutdown(ctx, func(context.Context) error {
		hookRan.Store(true)

		return nil
	})

	l.Go(func(context.Context) error {
		return nil
	})

	done := make(chan error, 1)
	go func() {
		done <- l.Run()
	}()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after heads finished")
	}

	require.True(t, hookRan.Load())
}

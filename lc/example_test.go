package lc_test

import (
	"context"
	"fmt"
	"time"

	"github.com/defany/goblin/lc"
)

func ExampleNew() {
	// Create a lifecycle with a custom shutdown timeout.
	app := lc.New(
		lc.WithShutdownTimeout(10 * time.Second),
	)

	// Register a goroutine that exits immediately.
	app.Go(func(_ context.Context) error {
		fmt.Println("worker started")
		return nil
	})

	err := app.Run()
	fmt.Println("error:", err)
	// Output:
	// worker started
	// error: <nil>
}

func ExampleLifecycle_Go_chaining() {
	// Chained heads run concurrently, but each waits for the previous to signal Ready or finish.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := lc.New(lc.WithContext(ctx))

	app.Go(func(ctx context.Context) error {
		fmt.Println("step 1: started")
		lc.Ready(ctx) // step 2 can now start
		return nil
	}).Go(func(ctx context.Context) error {
		fmt.Println("step 2: started")
		return nil
	})

	err := app.Run()
	fmt.Println("error:", err)
	// Output:
	// step 1: started
	// step 2: started
	// error: <nil>
}

func ExampleLifecycle_Go_chaining_withoutReady() {
	// Without Ready, the next head starts after the previous one finishes.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := lc.New(lc.WithContext(ctx))

	app.Go(func(_ context.Context) error {
		fmt.Println("step 1")
		return nil
	}).Go(func(_ context.Context) error {
		fmt.Println("step 2")
		return nil
	})

	err := app.Run()
	fmt.Println("error:", err)
	// Output:
	// step 1
	// step 2
	// error: <nil>
}

func ExampleDefer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := lc.New(lc.WithContext(ctx))

	app.Go(func(ctx context.Context) error {
		// Register a deferred cleanup using the lifecycle context.
		lc.Defer(app.Context(), func(_ context.Context) error {
			fmt.Println("cleanup executed")
			return nil
		})

		fmt.Println("main work")
		return nil
	})

	err := app.Run()
	fmt.Println("error:", err)
	// Output:
	// main work
	// cleanup executed
	// error: <nil>
}

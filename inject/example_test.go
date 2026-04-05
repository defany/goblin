package inject_test

import (
	"context"
	"fmt"

	"github.com/defany/goblin/inject"
)

func ExampleOnce() {
	ctx := context.Background()

	// Once ensures the factory runs only once per call site.
	// Subsequent calls from the same location return the cached value.
	val := inject.Once(ctx, func(_ context.Context) string {
		fmt.Println("factory called")
		return "singleton"
	})
	fmt.Println(val)
	// Output:
	// factory called
	// singleton
}

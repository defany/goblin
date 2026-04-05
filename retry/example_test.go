package retry_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/defany/goblin/retry"
)

func ExampleDo() {
	attempts := 0
	err := retry.Do(context.Background(), func(_ context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	},
		retry.WithAttempts(5),
		retry.WithBaseDelay(time.Millisecond),
		retry.WithoutJitter(),
	)
	fmt.Println("error:", err)
	fmt.Println("attempts:", attempts)
	// Output:
	// error: <nil>
	// attempts: 3
}

func ExampleDo_withRetryIf() {
	var errPermanent = errors.New("permanent")

	err := retry.Do(context.Background(), func(_ context.Context) error {
		return errPermanent
	},
		retry.WithAttempts(5),
		retry.WithBaseDelay(time.Millisecond),
		retry.WithRetryIf(func(err error) bool {
			// Only retry non-permanent errors.
			return !errors.Is(err, errPermanent)
		}),
	)
	fmt.Println(err)
	// Output: permanent
}

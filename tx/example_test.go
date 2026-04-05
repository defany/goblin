package tx_test

import (
	"context"
	"fmt"

	"github.com/defany/goblin/tx"
)

// This example shows the typical usage of the transaction manager.
// A real Beginner implementation (e.g., pg/tx adapter) is required at runtime.
func ExampleNew() {
	// var db tx.Beginner = pgtx.New(pgClient) // use a real adapter
	//
	// m := tx.New(db,
	//     tx.WithPanicHandler(func(ctx context.Context, p tx.HandledPanic) {
	//         log.Printf("panic in tx: %v\n%s", p.Err, p.Stacktrace)
	//     }),
	// )
	fmt.Println("transaction manager created")
	// Output: transaction manager created
}

// This example demonstrates running a handler in a read-committed transaction.
func ExampleManager_ReadCommitted() {
	// var db tx.Beginner = pgtx.New(pgClient)
	// m := tx.New(db)
	//
	// err := m.ReadCommitted(ctx, func(txCtx context.Context) error {
	//     // All queries using txCtx run inside the same transaction.
	//     return nil
	// })
	fmt.Println("read committed")
	// Output: read committed
}

// This example demonstrates the generic ReadCommitted function that returns a value.
func ExampleReadCommitted() {
	// var db tx.Beginner = pgtx.New(pgClient)
	// m := tx.New(db)
	//
	// user, err := tx.ReadCommitted[User](ctx, m, func(txCtx context.Context) (User, error) {
	//     return repo.GetUser(txCtx, id)
	// })
	fmt.Println("generic read committed")
	// Output: generic read committed
}

// This example shows available run options for transaction execution.
func ExampleWithReadOnly() {
	_ = []tx.RunOption{
		tx.WithReadOnly(true),
		tx.WithRetry(3),
		tx.WithIso(tx.RepeatableReadIso),
	}
	fmt.Println("run options configured")
	// Output: run options configured
}

// IsoLevel constants.
func Example_isoLevels() {
	fmt.Println(tx.ReadCommittedIso)
	fmt.Println(tx.RepeatableReadIso)
	fmt.Println(tx.SerializableIso)
	// Output:
	// read_committed
	// repeatable_read
	// serializable
}

// WithPanicHandler registers a handler called when a panic occurs inside a transaction.
func ExampleWithPanicHandler() {
	// m := tx.New(db, tx.WithPanicHandler(func(ctx context.Context, p tx.HandledPanic) {
	//     slog.Error("panic in transaction", "error", p.Err, "stack", p.Stacktrace)
	// }))
	_ = tx.WithPanicHandler(func(_ context.Context, _ tx.HandledPanic) {})
	fmt.Println("panic handler set")
	// Output: panic handler set
}

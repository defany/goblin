package errfmt_test

import (
	"errors"
	"fmt"

	"github.com/defany/goblin/errfmt"
)

func ExampleWithSource() {
	err := errors.New("connection refused")
	wrapped := errfmt.WithSource(err)
	// Result: errfmt/example_test.go:12 -> connection refused
	fmt.Println(wrapped)
}

func ExampleWithSource_withComment() {
	err := errors.New("not found")
	wrapped := errfmt.WithSource(err, "fetching user by ID")
	// Result: errfmt/example_test.go:19 -> [fetching user by ID] -> not found
	fmt.Println(wrapped)
}

func ExampleWithSource_nil() {
	wrapped := errfmt.WithSource(nil)
	fmt.Println(wrapped)
	// Output: <nil>
}

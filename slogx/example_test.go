package slogx_test

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/defany/goblin/slogx"
)

func ExampleJSON() {
	handler := slogx.JSON(
		slogx.WithLevel(slog.LevelInfo),
	)
	logger := slog.New(handler)
	// logger is ready to use with JSON output.
	_ = logger
	fmt.Println("json handler created")
	// Output: json handler created
}

func ExampleText() {
	handler := slogx.Text(
		slogx.WithLevel(slog.LevelWarn),
	)
	logger := slog.New(handler)
	// logger is ready to use with text output.
	_ = logger
	fmt.Println("text handler created")
	// Output: text handler created
}

func ExamplePretty() {
	handler := slogx.Pretty(
		slogx.WithLevel(slog.LevelDebug),
	)
	logger := slog.New(handler)
	// logger is ready to use with colorized pretty output.
	_ = logger
	fmt.Println("pretty handler created")
	// Output: pretty handler created
}

func ExampleErrAttr() {
	attr := slogx.ErrAttr(errors.New("something went wrong"))
	fmt.Println(attr.Key)
	fmt.Println(attr.Value.String())
	// Output:
	// error
	// something went wrong
}

func ExampleNewNoopLogger() {
	logger := slogx.NewNoopLogger()
	// The noop logger silently discards all log records.
	// Useful in tests or when logging must be disabled.
	logger.Info("this will be discarded")
	fmt.Println("noop logger created")
	// Output: noop logger created
}

func ExampleNewNoopHandler() {
	handler := slogx.NewNoopHandler()
	fmt.Println(handler.Enabled(nil, slog.LevelError))
	// Output: false
}

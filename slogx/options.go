package slogx

import (
	"io"
	"log/slog"
)

type Option func(*config)

type config struct {
	level     slog.Level
	addSource bool
	output    io.Writer
}

func WithLevel(l slog.Level) Option {
	return func(c *config) {
		c.level = l
	}
}

func WithAddSource(v bool) Option {
	return func(c *config) {
		c.addSource = v
	}
}

func WithOutput(w io.Writer) Option {
	return func(c *config) {
		c.output = w
	}
}

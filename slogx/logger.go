package slogx

import (
	"log/slog"
	"os"
)

func JSON(opts ...Option) slog.Handler {
	cfg := applyOpts(opts)
	return slog.NewJSONHandler(cfg.output, &slog.HandlerOptions{
		Level:     cfg.level,
		AddSource: cfg.addSource,
	})
}

func Text(opts ...Option) slog.Handler {
	cfg := applyOpts(opts)
	return slog.NewTextHandler(cfg.output, &slog.HandlerOptions{
		Level:     cfg.level,
		AddSource: cfg.addSource,
	})
}

func Pretty(opts ...Option) slog.Handler {
	cfg := applyOpts(opts)

	return NewPrettyHandler().
		WithOutput(cfg.output).
		WithLevel(cfg.level).
		WithAddSource(cfg.addSource).
		WithEmoji(true)
}

func defaults() config {
	return config{
		level:  slog.LevelDebug,
		output: os.Stdout,
	}
}

func applyOpts(opts []Option) config {
	cfg := defaults()
	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

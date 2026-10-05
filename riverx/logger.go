package riverx

import (
	"context"
	"log/slog"

	"github.com/defany/goblin/slogx"
	"github.com/riverqueue/river/riverlog"
)

func Logger(ctx context.Context) *slog.Logger {
	if logger, ok := riverlog.LoggerSafely(ctx); ok {
		return logger
	}

	return slogx.NewNoopLogger()
}

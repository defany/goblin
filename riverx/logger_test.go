package riverx

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger_NoopOutsideJob(t *testing.T) {
	logger := Logger(t.Context())

	require.NotNil(t, logger)
	require.False(t, logger.Enabled(t.Context(), slog.LevelError))
}

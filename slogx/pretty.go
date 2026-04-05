package slogx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	opts          slog.HandlerOptions
	logger        *log.Logger
	timeLayout    string
	attrs         []slog.Attr
	groups        []string
	useLevelEmoji bool
}

type levelInfo struct {
	text      string
	emoji     string
	colorFunc func(format string, a ...any) string
}

var levelsInfo = map[slog.Level]levelInfo{
	slog.LevelDebug: {"DEBUG", "👀", color.HiMagentaString},
	slog.LevelInfo:  {"INFO ", "💬", color.HiBlueString},
	slog.LevelWarn:  {"WARN ", "📛", color.HiYellowString},
	slog.LevelError: {"ERROR", "🔥", color.HiRedString},
}

func NewPrettyHandler() *PrettyHandler {
	return &PrettyHandler{
		logger:     log.New(os.Stderr, "", 0),
		opts:       slog.HandlerOptions{Level: slog.LevelDebug},
		timeLayout: time.RFC3339,
	}
}

func (h *PrettyHandler) WithOutput(w io.Writer) *PrettyHandler {
	h.logger = log.New(w, "", 0)
	return h
}

func (h *PrettyHandler) WithTimeLayout(layout string) *PrettyHandler {
	h.timeLayout = layout
	return h
}

func (h *PrettyHandler) WithLevel(l slog.Level) *PrettyHandler {
	h.opts.Level = l
	return h
}

func (h *PrettyHandler) WithAddSource(v bool) *PrettyHandler {
	h.opts.AddSource = v
	return h
}

func (h *PrettyHandler) WithEmoji(v bool) *PrettyHandler {
	h.useLevelEmoji = v

	return h
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	var parts []any

	if h.timeLayout != "" {
		parts = append(parts, color.WhiteString(r.Time.Format(h.timeLayout)))
	}

	parts = append(parts, h.formatLevel(r.Level), color.CyanString(r.Message))

	if attrs := h.formatAttrs(r); attrs != "" {
		parts = append(parts, attrs)
	}

	if h.opts.AddSource {
		parts = append(parts, color.GreenString(formatSource(r)))
	}

	h.logger.Println(parts...)
	return nil
}

func (h *PrettyHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l.Level() >= h.opts.Level.Level()
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := *h
	clone.attrs = append(clone.attrs, attrs...)
	return &clone
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	clone := *h
	clone.groups = append(clone.groups, name)
	return &clone
}

func (h *PrettyHandler) formatLevel(l slog.Level) string {
	info := levelsInfo[l]
	level := info.text
	if level == "" {
		level = l.String()
	}
	if info.colorFunc != nil {
		level = info.colorFunc(level)
	}
	if h.useLevelEmoji && info.emoji != "" {
		level = info.emoji + " " + level
	}
	return level
}

func (h *PrettyHandler) formatAttrs(r slog.Record) string {
	all := make([]slog.Attr, 0, r.NumAttrs()+len(h.attrs))
	r.Attrs(func(a slog.Attr) bool {
		all = append(all, a)
		return true
	})
	all = append(all, h.attrs...)

	fields := attrsToMap(all)
	if len(fields) == 0 {
		return ""
	}

	for i := len(h.groups) - 1; i >= 0; i-- {
		fields = map[string]any{h.groups[i]: fields}
	}

	b, err := json.Marshal(fields)
	if err != nil {
		return ""
	}
	return color.WhiteString(string(b))
}

func formatSource(r slog.Record) string {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	fn := filepath.Base(f.Function)
	if i := strings.IndexByte(fn, '.'); i >= 0 {
		fn = fn[i:]
	}

	return fmt.Sprintf("%s:%d%s", filepath.Base(f.File), f.Line, fn)
}

func attrsToMap(attrs []slog.Attr) map[string]any {
	fields := make(map[string]any, len(attrs))
	for _, a := range attrs {
		if a.Value.Kind() == slog.KindGroup {
			fields[a.Key] = attrsToMap(a.Value.Group())
		} else {
			fields[a.Key] = a.Value.Any()
		}
	}
	return fields
}

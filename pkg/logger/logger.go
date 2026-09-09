// Package logger provides a slog logger backed by zerolog.
package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

const _callerKey = "caller"

// New returns a slog logger that writes JSON to stdout.
func New(level string) *slog.Logger {
	return newWithWriter(os.Stdout, level)
}

func newWithWriter(w io.Writer, level string) *slog.Logger {
	zl := zerolog.New(w)

	handler := slogzerolog.Option{
		Level:     parseLevel(level),
		Logger:    &zl,
		AddSource: true,
		Converter: callerConverter,
	}.NewZerologHandler()

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// callerConverter replaces slog's "source" group with a flat "caller" field
// holding the base file name and line, matching the previous zerolog output.
func callerConverter(
	addSource bool,
	replaceAttr func(groups []string, a slog.Attr) slog.Attr,
	loggerAttr []slog.Attr,
	groups []string,
	record *slog.Record,
) map[string]any {
	out := slogzerolog.DefaultConverter(addSource, replaceAttr, loggerAttr, groups, record)

	if !addSource {
		return out
	}

	delete(out, slogzerolog.SourceKey)

	frames := runtime.CallersFrames([]uintptr{record.PC})
	frame, _ := frames.Next()
	out[_callerKey] = filepath.Base(frame.File) + ":" + strconv.Itoa(frame.Line)

	return out
}

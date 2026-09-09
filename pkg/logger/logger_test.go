package logger

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

var errTest = errors.New("test error")

func newBufferedLogger(level string) (*slog.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}

	return newWithWriter(buf, level), buf
}

func TestNew_WritesStructuredFields(t *testing.T) {
	t.Parallel()

	l, buf := newBufferedLogger("info")

	l.Info("hello", "user", "bob")

	out := buf.String()

	for _, want := range []string{
		`"level":"info"`,
		`"message":"hello"`,
		`"user":"bob"`,
		`"caller":"logger_test.go:`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output %q missing %q", out, want)
		}
	}
}

func TestNew_LevelFiltering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		level   string
		log     func(l *slog.Logger)
		wantLvl string
		want    bool
	}{
		{"debug suppressed at info", "info", func(l *slog.Logger) { l.Debug("dbg") }, `"level":"debug"`, false},
		{"debug emitted at debug", "debug", func(l *slog.Logger) { l.Debug("dbg") }, `"level":"debug"`, true},
		{"warn emitted at info", "info", func(l *slog.Logger) { l.Warn("warn") }, `"level":"warn"`, true},
		{"info suppressed at warn", "warn", func(l *slog.Logger) { l.Info("info") }, `"level":"info"`, false},
		{"error emitted at error", "error", func(l *slog.Logger) { l.Error("boom") }, `"level":"error"`, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l, buf := newBufferedLogger(tc.level)

			tc.log(l)

			if got := strings.Contains(buf.String(), tc.wantLvl); got != tc.want {
				t.Fatalf("level %q contains %s = %v, want %v (output: %s)",
					tc.level, tc.wantLvl, got, tc.want, buf.String())
			}
		})
	}
}

func TestNew_UnknownLevelDefaultsToInfo(t *testing.T) {
	t.Parallel()

	l, buf := newBufferedLogger("nonsense")

	l.Debug("dbg")
	l.Info("info")

	out := buf.String()

	if strings.Contains(out, `"level":"debug"`) {
		t.Fatalf("debug should be suppressed at default info level: %s", out)
	}

	if !strings.Contains(out, `"level":"info"`) {
		t.Fatalf("info should be emitted at default info level: %s", out)
	}
}

func TestNew_ErrorAttribute(t *testing.T) {
	t.Parallel()

	l, buf := newBufferedLogger("info")

	l.Error("amqp_rpc - V1 - createTask", "error", errTest)

	out := buf.String()

	if !strings.Contains(out, "test error") {
		t.Fatalf("error attribute missing from output: %s", out)
	}

	if !strings.Contains(out, `"level":"error"`) {
		t.Fatalf("error level missing from output: %s", out)
	}

	if strings.Contains(out, "%!(EXTRA") {
		t.Fatalf("output should not contain fmt leftovers: %s", out)
	}
}

func TestCallerConverter_WithoutSource(t *testing.T) {
	t.Parallel()

	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)

	out := callerConverter(false, nil, nil, nil, &rec)

	if _, ok := out[_callerKey]; ok {
		t.Fatalf("caller field should be omitted when source is disabled: %v", out)
	}
}

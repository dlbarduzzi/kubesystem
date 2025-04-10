package logging

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level string
	}{
		{
			name:  "dev",
			level: "debug",
		},
		{
			name:  "prod",
			level: "info",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := NewLogger(tt.name == "dev", tt.level)
			if logger == nil {
				t.Fatal("expected logger not to be nil")
			}
		})
	}
}

func TestDefaultLogger(t *testing.T) {
	t.Parallel()

	logger1 := DefaultLogger()
	if logger1 == nil {
		t.Fatal("expected logger not to be nil")
	}

	logger2 := DefaultLogger()
	if logger2 == nil {
		t.Fatal("expected logger not to be nil")
	}

	if logger1 != logger2 {
		t.Errorf("expected logger %#v to be equal %#v", logger1, logger2)
	}
}

func TestLoggerContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	logger1 := LoggerFromContext(ctx)
	if logger1 == nil {
		t.Fatal("expected logger not to be nil")
	}

	ctx = LoggerWithContext(ctx, logger1)

	logger2 := LoggerFromContext(ctx)
	if logger1 != logger2 {
		t.Errorf("expected logger %#v to be equal %#v", logger1, logger2)
	}
}

func TestReplaceAttr(t *testing.T) {
	t.Parallel()

	tm := time.Now().UTC()
	sr := &slog.Source{Function: "main.main", File: "/path/to/file", Line: 12}

	tests := []struct {
		name        string
		nano        bool
		loggerKey   string
		loggerValue slog.Value
		wantKey     string
		wantValue   string
	}{
		{
			name:        "time dev",
			loggerKey:   slog.TimeKey,
			loggerValue: slog.TimeValue(tm),
			wantValue:   tm.Format("2006-01-02T15:04:05"),
		},
		{
			name:        "time prod",
			nano:        true,
			loggerKey:   slog.TimeKey,
			loggerValue: slog.TimeValue(tm),
			wantValue:   tm.Format("2006-01-02T15:04:05.000000000Z"),
		},
		{
			name:      "message",
			loggerKey: slog.MessageKey,
			wantKey:   "message",
		},
		{
			name:        "source",
			loggerKey:   slog.SourceKey,
			loggerValue: slog.AnyValue(sr),
			wantKey:     "caller",
			wantValue:   "/path/to/file:12",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fn := replaceAttr(tt.nano)

			attr := slog.Attr{
				Key:   tt.loggerKey,
				Value: tt.loggerValue,
			}

			resp := fn(nil, attr)

			if tt.loggerKey == slog.TimeKey {
				got := resp.Value.String()
				if got != tt.wantValue {
					t.Errorf("expected time to be %s; got %s", tt.wantValue, got)
				}
			}

			if tt.loggerKey == slog.MessageKey {
				got := resp.Key
				if got != tt.wantKey {
					t.Errorf("expected message key to be %s; got %s", tt.wantKey, got)
				}
			}

			if tt.loggerKey == slog.SourceKey {
				gotKey := resp.Key
				if gotKey != tt.wantKey {
					t.Errorf("expected source key to be %s; got %s", tt.wantKey, gotKey)
				}
				gotValue := resp.Value
				if gotValue.String() != tt.wantValue {
					t.Errorf("expected source value to be %s; got %s", tt.wantValue, gotValue)
				}
			}
		})
	}
}

func TestGetLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		level     string
		wantLevel slog.Level
	}{
		{
			name:      "empty",
			level:     "",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "invalid",
			level:     "invalid",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "debug",
			level:     "debug",
			wantLevel: slog.LevelDebug,
		},
		{
			name:      "info",
			level:     "info",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "warn",
			level:     "warn",
			wantLevel: slog.LevelWarn,
		},
		{
			name:      "error",
			level:     "error",
			wantLevel: slog.LevelError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			level := getLogLevel(tt.level)
			if level != tt.wantLevel {
				t.Errorf("expected logger level to be %v; got %v", tt.wantLevel, level)
			}
		})
	}
}

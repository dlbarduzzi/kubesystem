package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

// contextKey is the logger string type used to avoid context collisions.
type contextKey string

// loggerKey identifies the logger value stored in the context.
const loggerKey = contextKey("logger")

var (
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

func NewLogger(dev bool, level string) *slog.Logger {
	var logger *slog.Logger

	if dev {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:       getLogLevel(level),
			AddSource:   true,
			ReplaceAttr: replaceAttr(false),
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       getLogLevel(level),
			AddSource:   true,
			ReplaceAttr: replaceAttr(true),
		}))
	}

	return logger
}

func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})
	return defaultLogger
}

func NewLoggerFromEnv() *slog.Logger {
	dev := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_MODE"))) == "dev"
	level := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	return NewLogger(dev, level)
}

func LoggerWithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return DefaultLogger()
}

type slogAttr func(groups []string, attr slog.Attr) slog.Attr

func replaceAttr(nano bool) slogAttr {
	return func(groups []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey {
			attr.Key = "time"
			attr.Value = slog.StringValue(getTimeFormat(attr.Value.Time().UTC(), nano))
		}
		if attr.Key == slog.MessageKey {
			attr.Key = "message"
		}
		if attr.Key == slog.SourceKey {
			source := attr.Value.Any().(*slog.Source)
			attr.Key = "caller"
			attr.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))
		}
		return attr
	}
}

const (
	levelDebug = "DEBUG"
	levelInfo  = "INFO"
	levelWarn  = "WARN"
	levelError = "ERROR"
)

func getLogLevel(level string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case levelDebug:
		return slog.LevelDebug
	case levelInfo:
		return slog.LevelInfo
	case levelWarn:
		return slog.LevelWarn
	case levelError:
		return slog.LevelError
	}
	return slog.LevelInfo
}

func getTimeFormat(t time.Time, nano bool) string {
	if nano {
		return t.Format("2006-01-02T15:04:05.000000000Z")
	} else {
		return t.Format("2006-01-02T15:04:05")
	}
}

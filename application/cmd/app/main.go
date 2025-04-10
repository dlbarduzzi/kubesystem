package main

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/dlbarduzzi/kubesystem/internal/logging"
)

func main() {
	logger := logging.NewLogger(false, "debug")

	ctx := context.Background()
	ctx = logging.LoggerWithContext(ctx, logger)

	cronjobGenerateLogs(ctx)
}

func cronjobGenerateLogs(ctx context.Context) {
	logger := logging.LoggerFromContext(ctx)
	for {
		ls := generateRandomLogStatus()
		status := slog.String("status", ls.status)

		user := slog.String("user", gofakeit.Email())
		message := gofakeit.Sentence(5)

		switch ls.level {
		case levelDebug:
			logger.Debug(message, status, user)
		case levelInfo:
			logger.Info(message, status, user)
		case levelWarn:
			logger.Warn(message, status, user)
		case levelError:
			logger.Error(message, status, user)
		default:
			logger.Info(message, status, user)
		}

		time.Sleep(time.Second * 1)
	}
}

const (
	levelDebug = "DEBUG"
	levelInfo  = "INFO"
	levelWarn  = "WARN"
	levelError = "ERROR"
)

var (
	debugStatuses = []string{"CACHE_HIT", "CONFIG_LOADED", "DB_CONNECTED"}
	infoStatuses  = []string{"USER_CREATED", "CHANNEL_REQUESTED", "PLAYLIST_STOPPED"}
	warnStatuses  = []string{"HIGH_MEMORY", "SLOW_RESPONSE", "DEPRECATED_API_CALL"}
	errorStatuses = []string{"TIMEOUT_OCCURRED", "UNHANDLED_EXCEPTION", "DB_CONNECTION_FAILED"}
)

type logStatus struct {
	level  string
	status string
}

func generateRandomLogStatus() logStatus {
	levels := []string{levelDebug, levelInfo, levelWarn, levelError}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(len(levels))

	level := levels[i]
	var status string

	switch level {
	case levelDebug:
		i = r.Intn(len(debugStatuses))
		status = debugStatuses[i]
	case levelInfo:
		i = r.Intn(len(infoStatuses))
		status = infoStatuses[i]
	case levelWarn:
		i = r.Intn(len(warnStatuses))
		status = warnStatuses[i]
	case levelError:
		i = r.Intn(len(errorStatuses))
		status = errorStatuses[i]
	default:
		i = r.Intn(len(infoStatuses))
		status = infoStatuses[i]
	}

	return logStatus{level, status}
}

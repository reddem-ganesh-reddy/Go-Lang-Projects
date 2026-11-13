package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	appLogger *slog.Logger
	initOnce  sync.Once
)

// Init creates the shared structured logger.
func Init(level string, serviceName ...string) error {
	var initErr error

	initOnce.Do(func() {
		var logLevel slog.Level

		switch strings.ToUpper(strings.TrimSpace(level)) {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "INFO":
			logLevel = slog.LevelInfo
		case "WARN":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		default:
			initErr = fmt.Errorf("invalid log level: %s", level)
			return
		}

		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})

		appLogger = slog.New(handler)

		if len(serviceName) > 0 && serviceName[0] != "" {
			appLogger = appLogger.With("service", serviceName[0])
		}
	})

	return initErr
}

// Get returns the shared logger.
func Get() *slog.Logger {
	if appLogger == nil {
		_ = Init("INFO")
	}

	return appLogger
}
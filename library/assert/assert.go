package assert

import (
	"log/slog"
	"os"
)

var assertLogger *slog.Logger = slog.Default()

func SetLogger(logger *slog.Logger) {
	assertLogger = logger
}

func Ok(err error, message string) {
	if err != nil {
		assertLogger.Error(message, slog.Any("err", err))
		os.Exit(1)
	}
}

func True(condition bool, message string) {
	if !condition {
		assertLogger.Error(message)
		os.Exit(1)
	}
}

func NotNil(object interface{}, message string) {
	if object == nil {
		assertLogger.Error(message)
		os.Exit(1)
	}
}

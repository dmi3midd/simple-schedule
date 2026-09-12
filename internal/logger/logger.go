package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Setup initializes the default slog logger with the specified log level.
// Supported levels: "debug", "info", "warn", "error". Defaults to "info".
func Setup(levelStr string) {
	var level slog.Level
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	slog.SetDefault(slog.New(handler))
}

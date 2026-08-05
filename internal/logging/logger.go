package logging

import (
	"log/slog"
	"os"
	"strings"
)

// Setup creates a new slog logger and sets it as the default.
func Setup(level string, isJSON bool) *slog.Logger {
	var slogLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn", "warning":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: slogLevel,
	}

	var handler slog.Handler
	if isJSON {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	// Wrap the handler with the redactor
	redactingHandler := NewRedactingHandler(handler)
	
	logger := slog.New(redactingHandler)
	slog.SetDefault(logger)

	return logger
}

package logging

import (
	"log/slog"
	"os"
)

// NewLogger creates a new slog.Logger configured for the given environment.
//
// For "production" it uses JSON output; otherwise it uses human-readable text.
func NewLogger(env string) *slog.Logger {
	var handler slog.Handler
	switch env {
	case "production":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}


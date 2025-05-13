package utils

import (
	// "fmt"
	"io"
	"log/slog"
)

// slog.LevelDebug, slog.LevelInfo
func NewJSONLogger(w io.Writer, level slog.Level) *slog.Logger {
	var levelVar *slog.LevelVar

	levelVar = new(slog.LevelVar)
	levelVar.Set(level)

	return slog.New(slog.NewJSONHandler(
		w,
		&slog.HandlerOptions{Level: levelVar},
	))
}

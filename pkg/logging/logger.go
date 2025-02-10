package logging

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func NewLogger() *slog.Logger {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)
	Logger = slog.New(handler)
	return Logger
}

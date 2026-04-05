package observability

import (
	"log"
	"log/slog"
	"os"
)

func NewLogger(env string) *slog.Logger {
	var level slog.Level

	if env == "development" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)

	return logger
}

func InitGlobalLogger(logger *slog.Logger) {
	slog.SetDefault(logger)
	log.SetOutput(os.Stdout)
}

package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func LogInit() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(file, opts)
	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

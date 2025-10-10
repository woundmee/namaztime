package slogger

import (
	"log/slog"
	"os"
)

// var Log *slog.Logger

func Init(logFilePath string) error {
	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) // 0644: -rw-r--r--
	if err != nil {
		return err
	}

	handler := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// сразу инициализирую, чтобы работать через Log
	slog.SetDefault(slog.New(handler))
	return nil
}

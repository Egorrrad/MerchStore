package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger

func Init(level string, format string, output string) {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var logOutput *os.File
	if output == "stdout" || output == "" {
		logOutput = os.Stdout
	} else {
		var err error
		logOutput, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic("Failed to open log file: " + err.Error())
		}
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(logOutput, &slog.HandlerOptions{Level: logLevel})
	default:
		handler = slog.NewTextHandler(logOutput, &slog.HandlerOptions{Level: logLevel})
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

package logger

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/truxcoder/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New() *log.Logger {
	wd, _ := os.Getwd()
	logFile := filepath.Join(wd, "temp", "app.log")

	hook := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	logger := log.New(log.DevelopmentMode, hook)
	logger.SetLevel(slog.LevelDebug)

	return logger
}

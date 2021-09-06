package log

import (
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

const (
	// InfoColor log
	InfoColor = "\033[1;34m%s\033[0m"
	// NoticeColor log
	NoticeColor = "\033[1;36m%s\033[0m"
	// WarningColor log
	WarningColor = "\033[1;33m%s\033[0m"
	// ErrorColor log
	ErrorColor = "\033[1;31m%s\033[0m"
	// DebugColor log
	DebugColor = "\033[0;36m%s\033[0m"
	// GreenColor log
	GreenColor = "\033[1;32m%s\033[0m"
)

// InitLogger init file
func InitLogger() {
	var logLevel = logrus.DebugLevel

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/console.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	logrus.SetLevel(logLevel)
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	logrus.AddHook(rotateFileHook)
}

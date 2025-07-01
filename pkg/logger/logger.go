package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes the logger with the given configuration
func InitLogger(level string, format string) error {
	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)

	// Set log format
	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}

	// Set output to stdout
	logrus.SetOutput(os.Stdout)

	return nil
}

// GetLogger returns a logger instance with context fields
func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

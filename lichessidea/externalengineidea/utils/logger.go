package utils

import (
	"go.uber.org/zap"
)

var (
	// Logger is the global logger instance
	Logger *zap.SugaredLogger
)

// InitLogger initializes the global logger
func InitLogger(development bool) error {
	var logger *zap.Logger
	var err error

	if development {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return err
	}

	Logger = logger.Sugar()
	return nil
}

// CloseLogger flushes any buffered log entries
func CloseLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

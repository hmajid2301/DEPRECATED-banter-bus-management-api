// Package core contains all the core high-level features.
package core

import (
	"io"

	"github.com/sirupsen/logrus"
)

// SetupLogger setups a new logger .
func SetupLogger(outputWriter io.Writer) (logger *logrus.Logger) {
	logger = logrus.New()
	logger.Out = outputWriter
	return logger
}

// UpdateLogLevel changes the log level for the logger.
func UpdateLogLevel(logger *logrus.Logger, logLevel string) *logrus.Logger {
	lLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		lLevel = logrus.InfoLevel
	}

	logger.SetLevel(lLevel)
	return logger
}

// UpdateFormatter changes the format for the logger depending on environment.
func UpdateFormatter(logger *logrus.Logger, environment string) *logrus.Logger {
	if environment == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	return logger
}

package core

import (
	"io"

	"github.com/sirupsen/logrus"
)

func SetupLogger(writer io.Writer) (logger *logrus.Logger) {
	logger = logrus.New()
	logger.Out = writer
	return logger
}

func UpdateLogLevel(logger *logrus.Logger, logLevel string) *logrus.Logger {
	lLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		lLevel = logrus.InfoLevel
	}

	logger.SetLevel(lLevel)
	return logger
}

func UpdateFormatter(logger *logrus.Logger, environment string) *logrus.Logger {
	if environment == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	return logger
}

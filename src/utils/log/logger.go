// Package logger is
package logger

import (
	"github.com/sirupsen/logrus"

	"banter-bus-server/src/utils/config"
)

// FormatLogger formats the logger as JSON/text depending on environment.
func FormatLogger(log *logrus.Logger) {
	appConfig := config.GetConfig()
	if appConfig.App.Environment == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}

	logLevel, err := logrus.ParseLevel(appConfig.App.LogLevel)
	if err != nil {
		logLevel = logrus.FatalLevel
	}

	log.SetLevel(logLevel)
}

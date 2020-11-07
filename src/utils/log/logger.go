// Package logger is
package logger

import (
	log "github.com/sirupsen/logrus"

	"banter-bus-server/src/utils/config"
)

// FormatLogger formats the logger as JSON/text depending on environment.
func FormatLogger() {
	appConfig := config.GetConfig()
	if appConfig.App.Environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	logLevel, err := log.ParseLevel(appConfig.App.LogLevel)
	if err != nil {
		logLevel = log.FatalLevel
	}

	log.SetLevel(logLevel)
}

// Package is ....
package logger

import (
	log "github.com/sirupsen/logrus"

	"banter-bus-server/src/utils/config"
)

func FormatLogger() {
	config := config.GetConfig()
	if config.App.Environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
	logLevel, err := log.ParseLevel(config.App.LogLevel)
	if err != nil {
		logLevel = log.FatalLevel
	}
	log.SetLevel(logLevel)
}

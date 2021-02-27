package controllers

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

// Env is related to all the data the controllers need.
type Env struct {
	Conf   core.Conf
	Logger *log.Logger
	DB     database.Database
}

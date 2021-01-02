package controllers

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
)

// Env is related to all the data the controllers need.
type Env struct {
	Config core.Config
	Logger *log.Logger
	DB     core.Repository
}

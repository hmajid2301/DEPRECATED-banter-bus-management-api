package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"

	"banter-bus-server/src/core/database"
	serverModels "banter-bus-server/src/server/models"
)

// Healthcheck checks if the API is healthy.
func Healthcheck(_ *gin.Context) (*serverModels.Healthcheck, error) {
	var err error

	healthy := database.Ping()

	if !healthy {
		return &serverModels.Healthcheck{}, errors.Errorf("Healthcheck Failed!")
	}

	return &serverModels.Healthcheck{
		Message: "The API is healthy.",
	}, err
}

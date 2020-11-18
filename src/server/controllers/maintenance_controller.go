package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server/models"
)

// Healthcheck checks if the API is healthy.
func Healthcheck(_ *gin.Context) (*models.Healthcheck, error) {
	var err error

	healthy := database.Ping()

	if !healthy {
		return &models.Healthcheck{}, errors.Errorf("Healthcheck Failed!")
	}

	return &models.Healthcheck{
		Message: "The API is healthy.",
	}, err
}

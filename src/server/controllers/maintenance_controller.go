package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server/models"
)

// Healthcheck checks if the API is healthy.
func Healthcheck(_ *gin.Context) (*models.Healthcheck, error) {
	var healthcheck = "The API is healthy."
	var err error

	healthy := database.Ping()

	if !healthy {
		healthcheck = "Database is not healthy!"
		err = errors.Errorf("Healthcheck Failed!")
	}
	return &models.Healthcheck{
		Message: healthcheck,
	}, err
}

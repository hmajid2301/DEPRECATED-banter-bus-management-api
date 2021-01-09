package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// Healthcheck checks if the API is healthy.
func (env *Env) Healthcheck(_ *gin.Context) (*serverModels.Healthcheck, error) {
	healthy := env.DB.Ping()

	if !healthy {
		return &serverModels.Healthcheck{}, errors.Errorf("Healthcheck Failed!")
	}

	return &serverModels.Healthcheck{
		Message: "The API is healthy.",
	}, nil
}

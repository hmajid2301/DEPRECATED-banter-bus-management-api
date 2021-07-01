package maintenance

import (
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

type MaintenanceAPI struct {
	Conf   core.Conf
	Logger *log.Logger
	DB     database.Database
}

func (env *MaintenanceAPI) Healthcheck(_ *gin.Context) (*Healthcheck, error) {
	healthy := env.DB.Ping()

	if !healthy {
		return &Healthcheck{}, errors.Errorf("Healthcheck Failed!")
	}

	return &Healthcheck{
		Message: "The API is healthy.",
	}, nil
}

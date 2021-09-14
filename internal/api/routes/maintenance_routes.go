package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/maintenance"
)

func MaintenanceRoutes(env *maintenance.MaintenanceAPI, grp *gin.RouterGroup) {
	grp.GET("/healthcheck", tonic.Handler(env.Healthcheck, http.StatusOK))
}

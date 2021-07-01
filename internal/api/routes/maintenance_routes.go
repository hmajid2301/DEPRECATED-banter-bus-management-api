package routes

import (
	"fmt"
	"net/http"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/maintenance"
)

// MaintenanceRoutes add routes related to the "maintenance" group.
func MaintenanceRoutes(env *maintenance.MaintenanceAPI, grp *fizz.RouterGroup) {
	grp.GET("/healthcheck", []fizz.OperationOption{
		fizz.Summary("Checks Banter Bus API is healthy."),
		fizz.Response(fmt.Sprint(http.StatusInternalServerError), "Server Error", nil, nil, nil),
	}, tonic.Handler(env.Healthcheck, http.StatusOK))
}

package routes

import (
	"fmt"
	"net/http"

	"banter-bus-server/src/server/controllers"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
)

// MaintenanceRoutes add routes related to the "maintenance" group.
func MaintenanceRoutes(grp *fizz.RouterGroup) {
	grp.GET("/healthcheck", []fizz.OperationOption{
		fizz.Summary("Checks Banter Bus API is healthy."),
		fizz.Response(fmt.Sprint(http.StatusInternalServerError), "Server Error", nil, nil),
	}, tonic.Handler(controllers.Healthcheck, http.StatusOK))
}

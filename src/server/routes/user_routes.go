package routes

import (
	"fmt"
	"net/http"

	"banter-bus-server/src/server/controllers"
	serverModels "banter-bus-server/src/server/models"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
)

// UserRoutes add routes related to the "user" group.
func UserRoutes(grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Add a new user."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil),
		fizz.Response(fmt.Sprint(http.StatusConflict), "User already exists", serverModels.APIError{}, nil),
	}, tonic.Handler(controllers.CreateUser, http.StatusCreated))
	grp.GET("", []fizz.OperationOption{
		fizz.Summary("Get all users."),
	}, tonic.Handler(controllers.GetAllUsers, http.StatusOK))
	grp.GET("/:name", []fizz.OperationOption{
		fizz.Summary("Get a user."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "User not found", serverModels.APIError{}, nil),
	}, tonic.Handler(controllers.GetUser, http.StatusOK))
}

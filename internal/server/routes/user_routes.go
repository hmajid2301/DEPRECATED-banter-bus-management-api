package routes

import (
	"fmt"
	"net/http"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/controllers"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// UserRoutes add routes related to the "user" group.
func UserRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Add a new user."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil, nil),
		fizz.Response(fmt.Sprint(http.StatusConflict), "User already exists", serverModels.APIError{}, nil, nil),
	}, tonic.Handler(env.CreateUser, http.StatusCreated))

	grp.GET("", []fizz.OperationOption{
		fizz.Summary("Get all users."),
	}, tonic.Handler(env.GetAllUsers, http.StatusOK))

	grp.GET("/:name", []fizz.OperationOption{
		fizz.Summary("Get a user."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "User not found", serverModels.APIError{}, nil, nil),
	}, tonic.Handler(env.GetUser, http.StatusOK))

	grp.DELETE("/:name", []fizz.OperationOption{
		fizz.Summary("Remove a user."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "User not found", serverModels.APIError{}, nil, nil),
	}, tonic.Handler(env.RemoveUser, http.StatusOK))
}

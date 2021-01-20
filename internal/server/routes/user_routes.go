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

	QuestionPoolRoutes(env, grp)

	grp.GET("/:name/story", []fizz.OperationOption{
		fizz.Summary("Get a user's stories."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "User not found", serverModels.APIError{}, nil, nil),
	}, tonic.Handler(env.GetUserStories, http.StatusOK))
}

// QuestionPoolRoutes add routes related question pools to the "user" group.
func QuestionPoolRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.GET("/:name/pool", []fizz.OperationOption{
		fizz.Summary("Get all of a user's question pools."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "User not found", serverModels.APIError{}, nil, nil),
	}, tonic.Handler(env.GetAllUserPools, http.StatusOK))

	grp.POST("/:name/pool", []fizz.OperationOption{
		fizz.Summary("Add a question pool to a user."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "User not found", serverModels.APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusConflict),
			"Question pool already exists",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.AddUserPool, http.StatusOK))

	grp.GET("/:name/pool/:pool_name", []fizz.OperationOption{
		fizz.Summary("Get a question pool from a user."),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"User or pool name not found",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.GetUserPool, http.StatusOK))

	grp.DELETE("/:name/pool/:pool_name", []fizz.OperationOption{
		fizz.Summary("Removes a question pool from a user."),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"User or pool name not found",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.RemoveUserPool, http.StatusOK))

	grp.PATCH("/:name/pool/:pool_name", []fizz.OperationOption{
		fizz.Summary("Add/remove questions in a question pool for a specific user."),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"User or pool name or question not found",
			serverModels.APIError{},
			nil,
			nil,
		),
		fizz.Response(
			fmt.Sprint(http.StatusConflict),
			"Question already exists",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.UpdateUserPool, http.StatusOK))
}

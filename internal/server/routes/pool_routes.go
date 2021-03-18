package routes

import (
	"fmt"
	"net/http"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/controllers"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// PoolRoutes add routes related question pools to the "user" group.
func PoolRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.GET("", []fizz.OperationOption{
		fizz.Summary("Get all of a user's pools."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "User not found", serverModels.APIError{}, nil, nil),
	}, tonic.Handler(env.GetAllPools, http.StatusOK))

	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Add a pool to a user."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "User not found", serverModels.APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusConflict),
			"Pool already exists",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.AddPool, http.StatusOK))

	grp.GET("/:pool_name", []fizz.OperationOption{
		fizz.Summary("Get a pool from a user."),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"User or pool name not found",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.GetPool, http.StatusOK))

	grp.DELETE("/:pool_name", []fizz.OperationOption{
		fizz.Summary("Removes a pool from a user."),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"User or pool name not found",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.RemovePool, http.StatusOK))

	PoolQuestionRoutes(env, grp)
}

// PoolQuestionRoutes are routes that can add or remove questions from a pool.
func PoolQuestionRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.POST("/:pool_name/question", []fizz.OperationOption{
		fizz.Summary("Add a question to the pool for a specific user."),
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
	}, tonic.Handler(env.AddQuestionToPool, http.StatusOK))

	grp.DELETE("/:pool_name/question", []fizz.OperationOption{
		fizz.Summary("Delete a question in a question pool for a specific user."),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"User or pool name or question not found",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.RemoveQuestionFromPool, http.StatusOK))
}

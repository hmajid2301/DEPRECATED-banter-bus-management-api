package routes

import (
	"fmt"
	"net/http"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/controllers"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// QuestionRoutes add routes related to the "question" group.
func QuestionRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Add a new question to a game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil, nil),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusConflict),
			"Question already exists for this game and this round",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.AddQuestion, http.StatusCreated))

	grp.GET("/group", []fizz.OperationOption{
		fizz.Summary("Get a list of question groups."),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or round does not exist",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.GetAllGroups, http.StatusOK))

	grp.DELETE("/:question_id", []fizz.OperationOption{
		fizz.Summary("Remove a question from a game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question doesn't exist",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.RemoveQuestion, http.StatusOK))

	translationRoutes(env, grp)
	updateRoutes(env, grp)
}

func translationRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.POST("/:question_id/:language", []fizz.OperationOption{
		fizz.Summary("Adds a new question translation."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question doesn't exist.",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.AddTranslation, http.StatusCreated))

	grp.DELETE("/:question_id/:language", []fizz.OperationOption{
		fizz.Summary("Remove a question translation from a game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question doesn't exist",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.RemoveTranslation, http.StatusOK))
}

func updateRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.PUT("/:question_id/enable", []fizz.OperationOption{
		fizz.Summary("Enables a question."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question does not exist",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.EnableQuestion, http.StatusOK))

	grp.PUT("/:question_id/disable", []fizz.OperationOption{
		fizz.Summary("Disabled a question."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question does not exist",
			serverModels.APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.DisableQuestion, http.StatusOK))
}

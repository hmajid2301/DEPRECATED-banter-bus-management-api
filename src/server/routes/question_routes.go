package routes

import (
	"fmt"
	"net/http"

	"banter-bus-server/src/server/controllers"
	serverModels "banter-bus-server/src/server/models"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
)

// QuestionRoutes add routes related to the "question" group.
func QuestionRoutes(grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Add a new question to a game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil),
		fizz.Response(
			fmt.Sprint(http.StatusConflict),
			"Question already exists for this game and this round",
			serverModels.APIError{},
			nil,
		),
	}, tonic.Handler(controllers.AddQuestion, http.StatusCreated))

	grp.DELETE("", []fizz.OperationOption{
		fizz.Summary("Remove a question from a game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question doesn't exist",
			serverModels.APIError{},
			nil,
		),
	}, tonic.Handler(controllers.RemoveQuestion, http.StatusOK))

	grp.PUT("", []fizz.OperationOption{
		fizz.Summary("Adds or removes a new question translation."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game doesn't exist or original question doesn't exist.",
			serverModels.APIError{},
			nil,
		),
	}, tonic.Handler(controllers.UpdateQuestion, http.StatusOK))

	grp.PUT("/enable", []fizz.OperationOption{
		fizz.Summary("Enables a question."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game or question do not exist", serverModels.APIError{}, nil),
	}, tonic.Handler(controllers.EnableQuestion, http.StatusOK))

	grp.PUT("/disable", []fizz.OperationOption{
		fizz.Summary("Disabled a question."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game or question do not exist", serverModels.APIError{}, nil),
	}, tonic.Handler(controllers.DisableQuestion, http.StatusOK))
}

package routes

import (
	"fmt"
	"net/http"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/questions"
)

func QuestionRoutes(env *questions.QuestionAPI, grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Add a new question to a game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", APIError{}, nil, nil),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusConflict),
			"Question already exists for this game and this round",
			APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.AddQuestion, http.StatusCreated))

	grp.GET("", []fizz.OperationOption{
		fizz.Summary("Gets a list of questions."),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound), "Game or round does not exist",
			APIError{}, nil, nil,
		),
	}, tonic.Handler(env.GetQuestions, http.StatusOK))

	grp.GET("/group", []fizz.OperationOption{
		fizz.Summary("Get a list of question groups."),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or round does not exist",
			APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.GetAllGroups, http.StatusOK))

	grp.GET("/:question_id/:language", []fizz.OperationOption{
		fizz.Summary("Get a single question."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game, question or language code (for that question) doesn't exist.",
			APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.GetQuestion, http.StatusOK))

	grp.DELETE("/:question_id", []fizz.OperationOption{
		fizz.Summary("Remove a question from a game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question doesn't exist",
			APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.RemoveQuestion, http.StatusOK))

	translationRoutes(env, grp)
	updateRoutes(env, grp)
}

func translationRoutes(env *questions.QuestionAPI, grp *fizz.RouterGroup) {
	grp.POST("/:question_id/:language", []fizz.OperationOption{
		fizz.Summary("Adds a new question translation."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question doesn't exist.",
			APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.AddTranslation, http.StatusCreated))

	grp.DELETE("/:question_id/:language", []fizz.OperationOption{
		fizz.Summary("Remove a question translation from a game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question doesn't exist",
			APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.RemoveTranslation, http.StatusOK))
}

func updateRoutes(env *questions.QuestionAPI, grp *fizz.RouterGroup) {
	grp.PUT("/:question_id/enable", []fizz.OperationOption{
		fizz.Summary("Enables a question."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question does not exist",
			APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.EnableQuestion, http.StatusOK))

	grp.PUT("/:question_id/disable", []fizz.OperationOption{
		fizz.Summary("Disabled a question."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", APIError{}, nil, nil),
		fizz.Response(
			fmt.Sprint(http.StatusNotFound),
			"Game or question does not exist",
			APIError{},
			nil,
			nil,
		),
	}, tonic.Handler(env.DisableQuestion, http.StatusOK))
}

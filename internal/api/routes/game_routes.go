package routes

import (
	"fmt"
	"net/http"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/games"
)

func GameRoutes(env *games.GameAPI, grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Create a new game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", APIError{}, nil, nil),
		fizz.Response(fmt.Sprint(http.StatusConflict), "Game already exists", APIError{}, nil, nil),
		fizz.Deprecated(true),
	}, tonic.Handler(env.AddGame, http.StatusCreated))

	grp.GET("", []fizz.OperationOption{
		fizz.Summary("Get all games."),
	}, tonic.Handler(env.GetGames, http.StatusOK))

	grp.GET("/:game_name", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", APIError{}, nil, nil),
		fizz.Summary("Get a game."),
	}, tonic.Handler(env.GetGame, http.StatusOK))

	grp.DELETE("/:game_name", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", APIError{}, nil, nil),
		fizz.Summary("Delete a game."),
		fizz.Deprecated(true),
	}, tonic.Handler(env.RemoveGame, http.StatusOK))

	grp.PUT("/:game_name/enable", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", APIError{}, nil, nil),
		fizz.Summary("Enables a game."),
	}, tonic.Handler(env.EnableGame, http.StatusOK))

	grp.PUT("/:game_name/disable", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", APIError{}, nil, nil),
		fizz.Summary("Disables a game."),
	}, tonic.Handler(env.DisableGame, http.StatusOK))
}

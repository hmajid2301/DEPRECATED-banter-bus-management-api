package routes

import (
	"fmt"
	"net/http"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/controllers"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// GameRoutes add routes related to the "game" group.
func GameRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Create a new game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil, nil),
		fizz.Response(fmt.Sprint(http.StatusConflict), "Game already exists", serverModels.APIError{}, nil, nil),
		fizz.Deprecated(true),
	}, tonic.Handler(env.CreateGame, http.StatusCreated))

	grp.GET("", []fizz.OperationOption{
		fizz.Summary("Get all games."),
	}, tonic.Handler(env.GetAllGames, http.StatusOK))

	grp.GET("/:name", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil, nil),
		fizz.Summary("Get a game."),
	}, tonic.Handler(env.GetGame, http.StatusOK))

	grp.DELETE("/:name", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil, nil),
		fizz.Summary("Delete a game."),
		fizz.Deprecated(true),
	}, tonic.Handler(env.RemoveGame, http.StatusOK))

	grp.PUT("/:name/enable", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil, nil),
		fizz.Summary("Enables a game."),
	}, tonic.Handler(env.EnableGame, http.StatusOK))

	grp.PUT("/:name/disable", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil, nil),
		fizz.Summary("Disables a game."),
	}, tonic.Handler(env.DisableGame, http.StatusOK))
}

package routes

import (
	"fmt"
	"net/http"

	"banter-bus-server/src/server/controllers"
	serverModels "banter-bus-server/src/server/models"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
)

// GameRoutes add routes related to the "game" group.
func GameRoutes(grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Create a new game."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", serverModels.APIError{}, nil),
		fizz.Response(fmt.Sprint(http.StatusConflict), "Game already exists", serverModels.APIError{}, nil),
		fizz.Deprecated(true),
	}, tonic.Handler(controllers.CreateGame, http.StatusCreated))

	grp.GET("", []fizz.OperationOption{
		fizz.Summary("Get all games."),
	}, tonic.Handler(controllers.GetAllGames, http.StatusOK))

	grp.GET("/:name", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil),
		fizz.Summary("Get a game."),
	}, tonic.Handler(controllers.GetGame, http.StatusOK))

	grp.DELETE("/:name", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil),
		fizz.Summary("Delete a game."),
		fizz.Deprecated(true),
	}, tonic.Handler(controllers.RemoveGame, http.StatusOK))

	grp.PUT("/:name/enable", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil),
		fizz.Summary("Enables a game."),
	}, tonic.Handler(controllers.EnableGame, http.StatusOK))

	grp.PUT("/:name/disable", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game doesn't exist", serverModels.APIError{}, nil),
		fizz.Summary("Disables a game."),
	}, tonic.Handler(controllers.DisableGame, http.StatusOK))
}

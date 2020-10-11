package server

import (
	"banter-bus-server/src/server/controllers"
	"banter-bus-server/src/server/models"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

func NewRouter() (*fizz.Fizz, error) {
	engine := gin.New()
	engine.Use(cors.Default())
	fizz := fizz.NewFromEngine(engine)

	infos := &openapi.Info{
		Title:       "Banter Bus",
		Description: "The API definition for the Banter Bus server.",
		Version:     "1.0.0",
	}
	fizz.GET("/openapi.json", nil, fizz.OpenAPI(infos, "json"))

	routes(fizz.Group("/game", "game", "Related to managing the game types."))

	if len(fizz.Errors()) != 0 {
		return nil, fmt.Errorf("fizz errors: %v", fizz.Errors())
	}
	return fizz, nil
}

func routes(grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Create a new game type."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad request", nil, nil),
		fizz.Response(fmt.Sprint(http.StatusConflict), "Does not exist", nil, nil),
	}, tonic.Handler(controllers.CreateGameType, http.StatusOK))
	tonic.SetErrorHook(errHook)
}

func errHook(c *gin.Context, e error) (int, interface{}) {
	code, msg := http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)

	if _, ok := e.(tonic.BindError); ok {
		code, msg = http.StatusBadRequest, e.Error()
	} else {
		switch {
		case errors.IsBadRequest(e), errors.IsNotValid(e), errors.IsNotSupported(e), errors.IsNotAssigned(e), errors.IsNotProvisioned(e):
			code, msg = http.StatusBadRequest, e.Error()
		case errors.IsForbidden(e):
			code, msg = http.StatusForbidden, e.Error()
		case errors.IsMethodNotAllowed(e):
			code, msg = http.StatusMethodNotAllowed, e.Error()
		case errors.IsNotFound(e), errors.IsUserNotFound(e):
			code, msg = http.StatusNotFound, e.Error()
		case errors.IsUnauthorized(e):
			code, msg = http.StatusUnauthorized, e.Error()
		case errors.IsAlreadyExists(e):
			code, msg = http.StatusConflict, e.Error()
		case errors.IsNotImplemented(e):
			code, msg = http.StatusNotImplemented, e.Error()
		}
	}
	err := models.ApiError{
		Message: msg,
	}
	return code, err
}

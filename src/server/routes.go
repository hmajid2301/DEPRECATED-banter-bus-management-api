package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"banter-bus-server/src/server/controllers"
	"banter-bus-server/src/server/models"
	logger "banter-bus-server/src/utils/log"
)

// NewRouter creates all the routes/endpoints, using Fizz.
func NewRouter() (*fizz.Fizz, error) {
	engine := gin.New()
	newLogger := logrus.New()
	logger.FormatLogger(newLogger)

	engine.Use(cors.Default())

	engine.Use(ginlogrus.Logger(newLogger), gin.Recovery())
	fizzApp := fizz.NewFromEngine(engine)

	infos := &openapi.Info{
		Title:       "Banter Bus",
		Description: "The API definition for the Banter Bus server.",
		Version:     "1.0.0",
	}
	fizzApp.GET("/openapi.json", nil, fizzApp.OpenAPI(infos, "json"))

	maintenanceRoutes(fizzApp.Group("", "maintenance", "Related to managing the maintenance of the API."))

	gameRoutes(fizzApp.Group("/game", "game", "Related to managing the game types."))
	questionRoutes(fizzApp.Group("/game", "question", "Related to managing the questions to a game type."))

	if len(fizzApp.Errors()) != 0 {
		return nil, fmt.Errorf("fizz errors: %v", fizzApp.Errors())
	}
	return fizzApp, nil
}

func maintenanceRoutes(grp *fizz.RouterGroup) {
	grp.GET("/healthcheck", []fizz.OperationOption{
		fizz.Summary("Checks Banter Bus API is healthy."),
		fizz.Response(fmt.Sprint(http.StatusInternalServerError), "Server Error", nil, nil),
	}, tonic.Handler(controllers.Healthcheck, http.StatusOK))
	tonic.SetErrorHook(errHook)
}

func gameRoutes(grp *fizz.RouterGroup) {
	grp.POST("", []fizz.OperationOption{
		fizz.Summary("Create a new game type."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", models.APIError{}, nil),
		fizz.Response(fmt.Sprint(http.StatusConflict), "Game type already exists", models.APIError{}, nil),
	}, tonic.Handler(controllers.CreateGameType, http.StatusCreated))
	grp.GET("", []fizz.OperationOption{
		fizz.Summary("Get all game types."),
	}, tonic.Handler(controllers.GetAllGameType, http.StatusOK))
	grp.GET("/:name", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game type doesn't exist", models.APIError{}, nil),
		fizz.Summary("Get a game types."),
	}, tonic.Handler(controllers.GetGameType, http.StatusOK))
	tonic.SetErrorHook(errHook)
	grp.DELETE("/:name", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game type doesn't exist", models.APIError{}, nil),
		fizz.Summary("Delete a game types."),
	}, tonic.Handler(controllers.RemoveGameType, http.StatusOK))
	grp.PUT("/:name/enable", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game type doesn't exist", models.APIError{}, nil),
		fizz.Response(fmt.Sprint(http.StatusConflict), "Game type is already enabled", models.APIError{}, nil),
		fizz.Summary("Enables a game type."),
	}, tonic.Handler(controllers.EnableGameType, http.StatusOK))
	grp.PUT("/:name/disable", []fizz.OperationOption{
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game type doesn't exist", models.APIError{}, nil),
		fizz.Response(fmt.Sprint(http.StatusConflict), "Game type is already disabled", models.APIError{}, nil),
		fizz.Summary("Disables a game type."),
	}, tonic.Handler(controllers.DisableGameType, http.StatusOK))
	tonic.SetErrorHook(errHook)
}

func questionRoutes(grp *fizz.RouterGroup) {
	grp.POST("/:name/question", []fizz.OperationOption{
		fizz.Summary("Add a new question to a game type."),
		fizz.Response(fmt.Sprint(http.StatusBadRequest), "Bad Request", models.APIError{}, nil),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Game type doesn't exist", models.APIError{}, nil),
		fizz.Response(
			fmt.Sprint(http.StatusConflict),
			"Question already exists for this game type and this round",
			models.APIError{},
			nil,
		),
	}, tonic.Handler(controllers.AddQuestion, http.StatusCreated))
	tonic.SetErrorHook(errHook)
}

func errHook(_ *gin.Context, e error) (int, interface{}) {
	code, msg := http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)

	if _, ok := e.(tonic.BindError); ok {
		code, msg = http.StatusBadRequest, e.Error()
	} else {
		switch {
		case errors.IsBadRequest(e), errors.IsNotValid(e), errors.IsNotSupported(e), errors.IsNotProvisioned(e):
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
	err := models.APIError{
		Message: msg,
	}
	return code, err
}

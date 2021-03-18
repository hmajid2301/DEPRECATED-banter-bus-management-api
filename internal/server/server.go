package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	ginlogrus "github.com/toorop/gin-logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/controllers"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/routes"
)

// Setup creates all the routes/endpoints, using Fizz.
func Setup(env *controllers.Env) (*fizz.Fizz, error) {
	engine := gin.New()

	if env.Conf.App.Env == "production" {
		engine.Use(gin.Recovery())
	}

	engine.Use(cors.Default())

	engine.Use(ginlogrus.Logger(env.Logger))
	fizzApp := fizz.NewFromEngine(engine)

	infos := &openapi.Info{
		Title:       "Banter Bus",
		Description: "The API definition for the Banter Bus server.",
		Version:     "1.0.0",
	}

	fizzApp.GET("/openapi.json", nil, fizzApp.OpenAPI(infos, "json"))
	routes.MaintenanceRoutes(env, fizzApp.Group("", "maintenance", "Related to managing the maintenance of the API."))
	routes.GameRoutes(env, fizzApp.Group("/game", "game", "Related to managing games."))
	routes.QuestionRoutes(env, fizzApp.Group("/game/:name/question", "question", "Related to managing the questions."))
	routes.UserRoutes(env, fizzApp.Group("/user", "user", "Related to managing users."))
	routes.StoryRoutes(env, fizzApp.Group("/user/:name/story", "user", "Related to managing stories."))
	routes.PoolRoutes(env, fizzApp.Group("/user/:name/pool", "user", "Related to managing question pools."))

	if len(fizzApp.Errors()) != 0 {
		return nil, fmt.Errorf("fizz errors: %v", fizzApp.Errors())
	}

	if env.Conf.App.Env != "production" {
		env.Logger.Info("activating pprof (devmode on)")
		pprof.Register(engine)
	}

	tonic.SetErrorHook(errHook)
	return fizzApp, nil
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
	err := serverModels.APIError{
		Message: msg,
	}
	return code, err
}

package api

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	ginlogrus "github.com/toorop/gin-logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/api/routes"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/maintenance"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/questions"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/story"
)

type Env struct {
	Conf   core.Conf
	Logger *log.Logger
	DB     database.Database
}

func Setup(env *Env) (*fizz.Fizz, error) {
	engine := gin.New()

	if env.Conf.App.Env == "production" {
		engine.Use(gin.Recovery())
	}

	engine.Use(cors.Default())

	engine.Use(ginlogrus.Logger(env.Logger))
	fizzApp := fizz.NewFromEngine(engine)

	infos := &openapi.Info{
		Title:       "Banter Bus Management API",
		Description: "The API specification for the Banter Bus Management API.",
		Version:     "1.0.0",
	}

	fizzApp.GET("/openapi", nil, fizzApp.OpenAPI(infos, "yaml"))
	routes.GameRoutes(&games.GameAPI{
		Conf:   env.Conf,
		Logger: env.Logger,
		DB:     env.DB,
	}, fizzApp.Group("/game", "game", "Related to managing games."))

	routes.QuestionRoutes(&questions.QuestionAPI{
		Conf:   env.Conf,
		Logger: env.Logger,
		DB:     env.DB,
	}, fizzApp.Group("/game/:game_name/question", "question", "Related to managing the questions."))

	routes.StoryRoutes(&story.StoryAPI{
		Conf:   env.Conf,
		Logger: env.Logger,
		DB:     env.DB,
	}, fizzApp.Group("/story/:game_name", "story", "Related to managing the stories."))

	routes.MaintenanceRoutes(&maintenance.MaintenanceAPI{
		Conf:   env.Conf,
		Logger: env.Logger,
		DB:     env.DB,
	}, fizzApp.Group("", "maintenance", "Related to maintenance of the app."))

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
	err := routes.APIError{
		Message: msg,
	}
	return code, err
}

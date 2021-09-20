package api

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	cors "github.com/rs/cors/wrapper/gin"
	ginlogrus "github.com/toorop/gin-logrus"

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

func Setup(env *Env) (*gin.Engine, error) {
	r := gin.New()

	if env.Conf.App.Env == "production" {
		r.Use(gin.Recovery())
	}

	corsConfig := cors.New(cors.Options{
		AllowedOrigins: env.Conf.Srv.CORS,
		Debug:          true,
	})

	r.Use(corsConfig)
	r.Use(ginlogrus.Logger(env.Logger))

	routes.GameRoutes(&games.GameAPI{
		Conf:   env.Conf,
		Logger: env.Logger,
		DB:     env.DB,
	}, r.Group("/game"))

	routes.QuestionRoutes(&questions.QuestionAPI{
		Conf:   env.Conf,
		Logger: env.Logger,
		DB:     env.DB,
	}, r.Group("/game/:game_name/question"))

	routes.StoryRoutes(&story.StoryAPI{
		Conf:   env.Conf,
		Logger: env.Logger,
		DB:     env.DB,
	}, r.Group("/story/:game_name"))

	routes.MaintenanceRoutes(&maintenance.MaintenanceAPI{
		Conf:   env.Conf,
		Logger: env.Logger,
		DB:     env.DB,
	}, r.Group(""))

	if env.Conf.App.Env != "production" {
		env.Logger.Info("activating pprof (devmode on)")
		pprof.Register(r)
	}

	tonic.SetErrorHook(errHook)
	return r, nil
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

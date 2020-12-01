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

	serverModels "banter-bus-server/src/server/models"
	"banter-bus-server/src/server/routes"
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
	routes.MaintenanceRoutes(fizzApp.Group("", "maintenance", "Related to managing the maintenance of the API."))
	routes.GameRoutes(fizzApp.Group("/game", "game", "Related to managing games."))
	routes.QuestionRoutes(fizzApp.Group("/game/:name/question", "question", "Related to managing the questions."))

	if len(fizzApp.Errors()) != 0 {
		return nil, fmt.Errorf("fizz errors: %v", fizzApp.Errors())
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

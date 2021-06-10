package routes

import (
	"fmt"
	"net/http"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/controllers"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// StoryRoutes add routes related to the "user" group.
func StoryRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.GET("", []fizz.OperationOption{
		fizz.Summary("Get a story."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Story not found", serverModels.APIError{}, nil, nil),
	}, tonic.Handler(env.GetStory, http.StatusOK))
}

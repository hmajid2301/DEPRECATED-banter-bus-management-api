package routes

import (
	"fmt"
	"net/http"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/controllers"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// StoryRoutes add routes related to the "story" group.
func StoryRoutes(env *controllers.Env, grp *fizz.RouterGroup) {
	grp.POST("/:name", []fizz.OperationOption{
		fizz.Summary("Add a story."),
	}, tonic.Handler(env.AddStory, http.StatusCreated))

	grp.GET("/:story_id", []fizz.OperationOption{
		fizz.Summary("Get a story."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Story not found", serverModels.APIError{}, nil, nil),
	}, tonic.Handler(env.GetStory, http.StatusOK))

	grp.DELETE("/:story_id", []fizz.OperationOption{
		fizz.Summary("Delete a story."),
		fizz.Response(fmt.Sprint(http.StatusNotFound), "Story not found", serverModels.APIError{}, nil, nil),
	}, tonic.Handler(env.DeleteStory, http.StatusOK))
}

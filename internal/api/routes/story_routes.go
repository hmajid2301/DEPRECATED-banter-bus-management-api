package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/story"
)

func StoryRoutes(env *story.StoryAPI, grp *gin.RouterGroup) {
	grp.POST("", tonic.Handler(env.AddStory, http.StatusCreated))
	grp.GET("/:story_id", tonic.Handler(env.GetStory, http.StatusOK))
	grp.DELETE("/:story_id", tonic.Handler(env.DeleteStory, http.StatusOK))
}

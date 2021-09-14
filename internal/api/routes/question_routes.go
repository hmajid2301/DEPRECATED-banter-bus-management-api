package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/questions"
)

func QuestionRoutes(env *questions.QuestionAPI, grp *gin.RouterGroup) {
	grp.POST("", tonic.Handler(env.AddQuestion, http.StatusCreated))
	grp.DELETE("/:question_id", tonic.Handler(env.RemoveQuestion, http.StatusOK))

	grp.GET("", tonic.Handler(env.GetQuestions, http.StatusOK))
	grp.GET("/group", tonic.Handler(env.GetAllGroups, http.StatusOK))
	grp.GET("/id", tonic.Handler(env.GetQuestionsIDs, http.StatusOK))
	grp.GET("/language", tonic.Handler(env.GetAllLanguages, http.StatusOK))

	grp.GET("/:question_id/:language", tonic.Handler(env.GetQuestion, http.StatusOK))
	grp.POST("/:question_id/:language", tonic.Handler(env.AddTranslation, http.StatusCreated))
	grp.DELETE("/:question_id/:language", tonic.Handler(env.RemoveTranslation, http.StatusOK))

	grp.PUT("/:question_id/enable", tonic.Handler(env.EnableQuestion, http.StatusOK))
	grp.PUT("/:question_id/disable", tonic.Handler(env.DisableQuestion, http.StatusOK))
}

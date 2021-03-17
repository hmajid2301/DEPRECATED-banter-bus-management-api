package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/factories"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// GetStories returns all the user's stories.
func (env *Env) GetStories(_ *gin.Context, params *serverModels.UserParams) ([]serverModels.Story, error) {
	userLogger := log.WithFields(log.Fields{
		"username": params.Username,
	})
	userLogger.Debug("Trying to get user stories.")

	s := service.StoryService{DB: env.DB, Username: params.Username}
	stories, err := s.Get()
	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("User does not exist."))
		return []serverModels.Story{}, errors.NotFoundf("the user %s", params.Username)
	}

	srvStories, err := env.newStories(stories)
	return srvStories, err
}

func (env *Env) newStories(userStories []models.Story) ([]serverModels.Story, error) {
	stories := []serverModels.Story{}

	for _, story := range userStories {
		game, err := factories.GetGame(story.GameName)
		if err != nil {
			env.Logger.Errorf("unknown game %s", story.GameName)
			return []serverModels.Story{}, err
		}

		newStory, err := game.NewStory(story)
		if err != nil {
			return []serverModels.Story{}, err
		}
		stories = append(stories, newStory)
	}

	return stories, nil
}

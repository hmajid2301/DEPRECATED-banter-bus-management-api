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

// GetStory returns the story with a specific ID.
func (env *Env) GetStory(_ *gin.Context, params *serverModels.StoryIDParams) (serverModels.Story, error) {
	storyLogger := log.WithFields(log.Fields{
		"story_id": params.StoryID,
	})
	storyLogger.Debug("Trying to get story.")

	s := service.StoryService{DB: env.DB}
	story, err := s.Get(params.StoryID)
	if err != nil {
		storyLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("Story does not exist."))
		return serverModels.Story{}, errors.NotFoundf("the story %s", params.StoryID)
	}

	srvstory, err := env.newStory(story)
	return srvstory, err
}

func (env *Env) newStory(story models.Story) (serverModels.Story, error) {
	game, err := factories.GetGame(story.GameName)
	if err != nil {
		env.Logger.Errorf("unknown game %s", story.GameName)
		return serverModels.Story{}, err
	}

	newStory, err := game.NewStory(story)
	if err != nil {
		return serverModels.Story{}, err
	}

	return newStory, nil
}

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

// AddStory adds a story.
func (env *Env) AddStory(_ *gin.Context, input *serverModels.NewStoryInput) (string, error) {
	story := input.Story
	gameName := input.GameParams.Name

	storyLogger := log.WithFields(log.Fields{
		"game_name": gameName,
		"story":     story,
	})
	storyLogger.Debug("Trying to add story.")

	s := service.StoryService{DB: env.DB}
	serviceStory, err := env.newStory(gameName, story)
	if err != nil {
		return "", err
	}

	id, err := s.Add(serviceStory)
	if err != nil {
		storyLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("Failed to add story."))
	}

	return id, nil
}

func (env *Env) newStory(gameName string, story serverModels.Story) (models.Story, error) {
	game, err := factories.GetGame(gameName)
	if err != nil {
		env.Logger.Errorf("unknown game %s", gameName)
		return models.Story{}, err
	}

	newStory, err := game.NewStory(story)
	if err != nil {
		return models.Story{}, err
	}

	return newStory, nil
}

// GetStory returns the story with a specific ID.
func (env *Env) GetStory(_ *gin.Context, story *serverModels.StoryIDParams) (serverModels.Story, error) {
	storyLogger := log.WithFields(log.Fields{
		"story_id": story.StoryID,
	})
	storyLogger.Debug("Trying to get story.")

	s := service.StoryService{DB: env.DB}
	stories, err := s.Get(story.StoryID)
	if err != nil {
		storyLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("Story does not exist."))
		return serverModels.Story{}, errors.NotFoundf("the story %s", story.StoryID)
	}

	srvStories, err := env.newServerStory(stories)
	return srvStories, err
}

func (env *Env) newServerStory(story models.Story) (serverModels.Story, error) {
	game, err := factories.GetGame(story.GameName)
	if err != nil {
		env.Logger.Errorf("unknown game %s", story.GameName)
		return serverModels.Story{}, err
	}

	newStory, err := game.NewServerStory(story)
	if err != nil {
		return serverModels.Story{}, err
	}

	return newStory, nil
}

// DeleteStory removed a story with a specific ID.
func (env *Env) DeleteStory(_ *gin.Context, params *serverModels.StoryIDParams) error {
	storyLogger := log.WithFields(log.Fields{
		"story_id": params.StoryID,
	})
	storyLogger.Debug("Trying to remove story.")

	s := service.StoryService{DB: env.DB}
	err := s.Delete(params.StoryID)
	if err != nil {
		storyLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("Story does not exist."))
		return errors.NotFoundf("the story with id %s", params.StoryID)
	}

	return nil
}

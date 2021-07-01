package story

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

// StoryAPI is related to all the data the controllers need.
type StoryAPI struct {
	Conf   core.Conf
	Logger *log.Logger
	DB     database.Database
}

func (env *StoryAPI) AddStory(_ *gin.Context, input *NewStoryInput) (string, error) {
	story := input.StoryInOut
	gameName := input.GameParams.Name

	storyLogger := env.Logger.WithFields(log.Fields{
		"game_name": gameName,
		"story":     story,
	})
	storyLogger.Debug("Trying to add story.")

	s := StoryService{DB: env.DB}
	serviceStory, err := env.newStory(gameName, story)
	if err != nil {
		storyLogger.WithFields(log.Fields{
			"err": err,
		}).Error(("Failed to convert story."))
		return "", err
	}

	id, err := s.Add(serviceStory)
	if err != nil {
		storyLogger.WithFields(log.Fields{
			"err": err,
		}).Error(("Failed to add story."))
	}

	return id, nil
}

func (env *StoryAPI) newStory(gameName string, story StoryInOut) (Story, error) {
	game, err := GetGame(gameName)
	if err != nil {
		return Story{}, errors.Errorf("unknown game %s, %v", gameName, err)
	}

	newStory, err := game.NewStory(story)
	if err != nil {
		return Story{}, err
	}

	return newStory, nil
}

func (env *StoryAPI) GetStory(_ *gin.Context, story *StoryIDParams) (StoryInOut, error) {
	storyLogger := env.Logger.WithFields(log.Fields{
		"story_id": story.StoryID,
	})
	storyLogger.Debug("Trying to get story.")

	s := StoryService{DB: env.DB}
	stories, err := s.Get(story.StoryID)
	if err != nil {
		storyLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("Story does not exist."))
		return StoryInOut{}, errors.NotFoundf("the story %s", story.StoryID)
	}

	srvStories, err := env.newAPIStory(stories)
	if err != nil {
		storyLogger.Errorf("Failed to convert Story %v", err)
		return StoryInOut{}, err
	}
	return srvStories, err
}

func (env *StoryAPI) newAPIStory(story Story) (StoryInOut, error) {
	game, err := GetGame(story.GameName)
	if err != nil {
		return StoryInOut{}, errors.Errorf("unknown game %s, %v", story.GameName, err)
	}

	newStory, err := game.NewStoryOut(story)
	if err != nil {
		return StoryInOut{}, err
	}

	return newStory, nil
}

func (env *StoryAPI) DeleteStory(_ *gin.Context, params *StoryIDParams) error {
	storyLogger := env.Logger.WithFields(log.Fields{
		"story_id": params.StoryID,
	})
	storyLogger.Debug("Trying to remove story.")

	s := StoryService{DB: env.DB}
	err := s.Delete(params.StoryID)
	if err != nil {
		storyLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("Story does not exist."))
		return errors.NotFoundf("the story with id %s", params.StoryID)
	}

	return nil
}

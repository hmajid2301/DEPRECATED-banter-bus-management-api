package story

import (
	"strings"

	"github.com/google/uuid"
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

type StoryService struct {
	DB database.Database
}

func (s *StoryService) Add(story Story) (string, error) {
	uuidWithHyphen := uuid.New()
	uuid := strings.ReplaceAll(uuidWithHyphen.String(), "-", "")
	story.ID = uuid

	inserted, err := story.Add(s.DB)
	if !inserted || err != nil {
		return "", errors.Errorf("failed to add story %v", err)
	}

	return uuid, nil
}

func (s *StoryService) Get(storyID string, gameName string) (Story, error) {
	filter := map[string]interface{}{
		"id":        storyID,
		"game_name": gameName,
	}

	story := Story{}
	err := story.Get(s.DB, filter)
	if err != nil {
		return Story{}, err
	}

	return story, nil
}

func (s *StoryService) Delete(storyID string, gameName string) error {
	filter := map[string]interface{}{
		"id":        storyID,
		"game_name": gameName,
	}

	deleted, err := s.DB.Delete("story", filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove story %v", err)
	}

	return nil
}

type Gamer interface {
	NewStory(story StoryInOut) (Story, error)
	NewStoryOut(story Story) (StoryInOut, error)
}

func GetGame(name string) (Gamer, error) {
	switch name {
	case "quibly":
		return Quibly{}, nil
	case "fibbing_it":
		return FibbingIt{}, nil
	case "drawlosseum":
		return Drawlosseum{}, nil
	default:
		return nil, errors.NotFoundf("Game %s", name)
	}
}

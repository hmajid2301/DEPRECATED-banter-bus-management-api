package story

import (
	"strings"

	"github.com/google/uuid"
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

// StoryService is struct data required by all story service functions.
type StoryService struct {
	DB database.Database
}

// Add is used to add a new story.
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

// Get get a specific story using id.
func (s *StoryService) Get(storyID string) (Story, error) {
	filter := map[string]string{
		"id": storyID,
	}

	story := Story{}
	err := story.Get(s.DB, filter)
	if err != nil {
		return Story{}, err
	}

	return story, nil
}

// Delete removes a specific story using id.
func (s *StoryService) Delete(storyID string) error {
	filter := map[string]string{
		"id": storyID,
	}

	deleted, err := s.DB.Delete("story", filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove story %v", err)
	}

	return nil
}

// Gamer is the interface for game(s).
type Gamer interface {
	NewStory(story StoryInOut) (Story, error)
	NewStoryOut(story Story) (StoryInOut, error)
}

// GetGame is the factory function which will return the game struct based on the name.
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

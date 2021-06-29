package service

import (
	"github.com/juju/errors"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// StoryService is struct data required by all story service functions.
type StoryService struct {
	DB database.Database
}

// Get get a specific story using id.
func (s *StoryService) Get(storyID string) (models.Story, error) {
	filter := map[string]string{
		"id": storyID,
	}

	story := models.Story{}
	err := story.Get(s.DB, filter)
	if err != nil {
		return models.Story{}, err
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

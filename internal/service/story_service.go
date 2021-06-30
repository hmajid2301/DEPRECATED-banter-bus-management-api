package service

import (
	"strings"

	"github.com/google/uuid"
	"github.com/juju/errors"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// StoryService is struct data required by all story service functions.
type StoryService struct {
	DB database.Database
}

// Add is used to add a new story.
func (s *StoryService) Add(story models.Story) (string, error) {
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

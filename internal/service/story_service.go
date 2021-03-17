package service

import (
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// StoryService is struct data required by all story service functions.
type StoryService struct {
	DB       database.Database
	Username string
}

// Get gets a specific user's stories.
func (s *StoryService) Get() (models.Stories, error) {
	u := UserService{DB: s.DB, Username: s.Username}
	_, err := u.Get()
	if err != nil {
		return nil, err
	}

	filter := map[string]string{
		"username": s.Username,
	}
	stories := models.Stories{}
	err = stories.Get(s.DB, filter)
	if err != nil {
		return models.Stories{}, err
	}

	return stories, nil
}

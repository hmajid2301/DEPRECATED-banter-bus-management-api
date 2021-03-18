package service

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// UserService is struct data required by all user service functions.
type UserService struct {
	DB       database.Database
	Username string
}

// Add is used to add a new user type
func (u *UserService) Add(membership string, admin *bool) error {
	exists, _ := u.doesUserExist()
	if exists {
		return errors.AlreadyExistsf("the user %s", u.Username)
	}

	var user = models.User{
		Username:    u.Username,
		Admin:       admin,
		Privacy:     "private",
		Membership:  membership,
		Preferences: &models.UserPreferences{},
		Friends:     []models.Friend{},
	}

	inserted, err := user.Add(u.DB)
	if !inserted {
		return errors.Errorf("failed to add user %s", u.Username)
	}
	return err
}

// Remove removes a user from the database
func (u *UserService) Remove() error {
	user, err := u.Get()
	if user == nil {
		return errors.NotFoundf("the user %s", u.Username)
	} else if err != nil {
		return err
	}

	for _, pool := range user.QuestionPools {
		p := PoolService{Username: pool.Username, PoolName: pool.PoolName}
		err := p.RemovePool()

		if err != nil {
			return errors.Errorf("failed to remove pool %s from user %s", pool.PoolName, u.Username)
		}
	}

	filter := map[string]string{"username": u.Username}
	deleted, err := u.DB.Delete("user", filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove user %s", u.Username)
	}

	return nil
}

// GetAll is used to get usernames of all users
func (u *UserService) GetAll(adminFilter *bool, privacyFilter *string, membershipFilter *string) ([]string, error) {
	users := models.Users{}

	emptyFilter := map[string]string{}
	err := users.Get(u.DB, emptyFilter)
	if err != nil {
		return []string{}, err
	}

	var usernames []string
	for _, user := range users {
		if (adminFilter == nil || *adminFilter == *user.Admin) &&
			(privacyFilter == nil || *privacyFilter == user.Privacy) &&
			(membershipFilter == nil || *membershipFilter == user.Membership) {
			usernames = append(usernames, user.Username)
		}
	}
	return usernames, nil
}

// Get gets a user with a given username
func (u *UserService) Get() (*models.User, error) {
	var (
		filter = map[string]string{"username": u.Username}
		user   = &models.User{}
	)

	err := user.Get(u.DB, filter)
	if err != nil {
		return &models.User{}, errors.NotFoundf("user %s, %s", u.Username, err)
	}

	return user, nil
}

func (u *UserService) doesUserExist() (bool, error) {
	user, err := u.Get()
	exists := (user.Username != "")
	return exists, err
}

func validateQuestion(gameName string, question models.GenericQuestion) error {
	game, err := games.GetGame(gameName)
	if err != nil {
		return err
	}

	err = game.ValidateQuestion(question)
	if err != nil {
		return err
	}

	return nil
}

package biz

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
)

// UserService is struct data required by all user service functions.
type UserService struct {
	DB core.Repository
}

// Add is used to add a new user type
func (u *UserService) Add(username string, membership string, admin *bool) error {
	if u.doesUserExist(username) {
		return errors.AlreadyExistsf("The user %s", username)
	}
	var newUser = models.User{
		Username:      username,
		Admin:         admin,
		Privacy:       "public",
		Membership:    membership,
		Preferences:   &models.UserPreferences{},
		Friends:       []models.Friend{},
		Stories:       []models.Story{},
		QuestionPools: []models.QuestionPool{},
	}

	inserted, err := u.DB.Insert("user", newUser)
	if !inserted {
		logger := log.WithFields(log.Fields{
			"username": username,
			"err":      err,
		})
		logger.Error("Error:")
		return errors.Errorf("Failed to add new user %s", username)
	}
	return err
}

// GetAll is used to get usernames of all users
func (u *UserService) GetAll(adminFilter *bool, privacyFilter *string, membershipFilter *string) ([]string, error) {
	users := []models.User{}

	err := u.DB.GetAll("user", &users)
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
func (u *UserService) Get(username string) (*models.User, error) {
	var (
		filter = map[string]string{"username": username}
		user   *models.User
	)

	err := u.DB.Get("user", filter, &user)
	if err != nil {
		return &models.User{}, errors.NotFoundf("User %s", username)
	}

	return user, nil
}

// GetPools gets a specific user's question pools
func (u *UserService) GetPools(username string) ([]models.QuestionPool, error) {
	user, err := u.Get(username)
	if err != nil {
		return []models.QuestionPool{}, err
	}
	pools := user.QuestionPools
	for index, pool := range pools {
		genericQuestions, err := u.newGenericQuestions(pool.GameName, pool.Questions)
		if err != nil {
			return []models.QuestionPool{}, err
		}
		pools[index].Questions = genericQuestions
	}

	return pools, nil
}

// Remove removes a user from the database
func (u *UserService) Remove(username string) error {
	if !u.doesUserExist(username) {
		return errors.NotFoundf("The user %s", username)
	}

	filter := map[string]string{"username": username}

	deleted, err := u.DB.Delete("user", filter)
	if !deleted || err != nil {
		return errors.Errorf("Failed to remove user %s", username)
	}

	return nil
}

func (u *UserService) newGenericQuestions(name string, questions interface{}) ([]models.GenericQuestion, error) {
	var genericQuestions []models.GenericQuestion
	gameType, err := getGameType(name, models.GenericQuestion{}, u.DB)
	if err != nil {
		return nil, err
	}

	genericQuestions, err = gameType.QuestionPoolToGenericQuestions(questions)
	return genericQuestions, err
}

func (u *UserService) doesUserExist(username string) bool {
	user, _ := u.Get(username)
	return user.Username != ""
}

// GetUserStories gets a specific user's stories.
func (u *UserService) GetUserStories(username string) ([]models.Story, error) {
	user, err := u.Get(username)
	if err != nil {
		return []models.Story{}, err
	}

	return user.Stories, nil
}

package core

import (
	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/models"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

// AddUser is used to add a new user type
func AddUser(username string, membership string, admin *bool) error {
	if doesUserExist(username) {
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

	inserted, err := database.Insert("user", newUser)
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

// GetAllUsers is used to get usernames of all users
func GetAllUsers(adminFilter *bool, privacyFilter *string, membershipFilter *string) ([]string, error) {
	users := []models.User{}

	err := database.GetAll("user", &users)
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

// GetUser gets a user with a given username
func GetUser(username string) (*models.User, error) {
	var (
		filter = map[string]string{"username": username}
		user   *models.User
	)

	err := database.Get("user", filter, &user)
	if err != nil {
		return &models.User{}, errors.NotFoundf("User %s", username)
	}

	return user, nil
}

// GetUserPools gets a specific user's question pools
func GetUserPools(username string) ([]models.QuestionPool, error) {
	user, err := GetUser(username)
	if err != nil {
		return []models.QuestionPool{}, err
	}
	pools := user.QuestionPools
	for index, pool := range pools {
		genericQuestions, err := newGenericQuestions(pool.GameName, pool.Questions)
		if err != nil {
			return []models.QuestionPool{}, err
		}
		pools[index].Questions = genericQuestions
	}

	return pools, nil
}

// RemoveUser removes a user from the database
func RemoveUser(username string) error {
	if !doesUserExist(username) {
		return errors.NotFoundf("The user %s", username)
	}

	filter := map[string]string{"username": username}

	deleted, err := database.Delete("user", filter)
	if !deleted || err != nil {
		return errors.Errorf("Failed to remove user %s", username)
	}

	return nil
}

func newGenericQuestions(name string, questions interface{}) ([]models.GenericQuestion, error) {
	var genericQuestions []models.GenericQuestion
	gameType, err := getGameType(name, models.GenericQuestion{})
	if err != nil {
		return nil, err
	}

	genericQuestions, err = gameType.QuestionPoolToGenericQuestions(questions)
	return genericQuestions, err
}

func doesUserExist(username string) bool {
	user, _ := GetUser(username)
	return user.Username != ""
}

// GetUserStories gets a specific user's stories.
func GetUserStories(username string) ([]models.Story, error) {
	user, err := GetUser(username)
	if err != nil {
		return []models.Story{}, err
	}

	return user.Stories, nil
}

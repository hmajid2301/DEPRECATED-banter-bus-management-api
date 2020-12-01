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
		filter = models.User{Username: username}
		user   *models.User
	)

	err := database.Get("user", filter, &user)
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func doesUserExist(username string) bool {
	user, _ := GetUser(username)
	return user.Username != ""
}

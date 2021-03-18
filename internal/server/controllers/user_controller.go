package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// CreateUser adds a new user to the database
func (env *Env) CreateUser(_ *gin.Context, user *serverModels.NewUser) error {
	userLogger := log.WithFields(log.Fields{
		"username": user.Username,
	})
	userLogger.Debug("Trying to add new user.")
	u := service.UserService{DB: env.DB, Username: user.Username}
	err := u.Add(user.Membership, user.Admin)

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to add new user.")

		if errors.IsAlreadyExists(err) {
			userLogger.WithFields(log.Fields{
				"err": err,
			}).Warn("User already exists.")
		}
		return err
	}

	return nil
}

// GetAllUsers gets a list of all usernames.
func (env *Env) GetAllUsers(_ *gin.Context, params *serverModels.ListUserParams) ([]string, error) {
	log.Debug("Trying to get all users")
	var (
		t = true
		f = false
	)

	var n *bool

	// This filter for converting admin requests into bool
	adminFilters := map[string]*bool{
		"admin":     &t,
		"non_admin": &f,
		"all":       n,
	}

	privacyFilter := newUserFilter(params.Privacy)
	membershipFilter := newUserFilter(params.Membership)
	adminFilter := adminFilters[params.AdminStatus]

	u := service.UserService{DB: env.DB}
	usernames, err := u.GetAll(adminFilter, privacyFilter, membershipFilter)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get user.")
		return []string{}, err
	}

	return usernames, nil
}

func newUserFilter(param string) *string {
	var newFilter *string
	if param != "all" {
		newFilter = &param
	}
	return newFilter
}

// GetUser gets all the information about a specific user.
func (env *Env) GetUser(_ *gin.Context, params *serverModels.UserParams) (*serverModels.User, error) {
	userLogger := log.WithFields(log.Fields{
		"username": params.Username,
	})
	userLogger.Debug("Trying to add new user.")

	u := service.UserService{DB: env.DB, Username: params.Username}
	user, err := u.Get()
	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("User does not exist."))
		return &serverModels.User{}, errors.NotFoundf("The user %s", params.Username)
	}

	srvUser := newUser(user)
	return srvUser, nil
}

// RemoveUser removes a user
func (env *Env) RemoveUser(_ *gin.Context, user *serverModels.UserParams) error {
	userLogger := log.WithFields(log.Fields{
		"username": user.Username,
	})
	userLogger.Debug("Trying to add new user.")

	userService := service.UserService{DB: env.DB, Username: user.Username}
	err := userService.Remove()
	return err
}

func newUser(userFromDB *models.User) *serverModels.User {
	returnedUser := &serverModels.User{
		Username:   userFromDB.Username,
		Admin:      userFromDB.Admin,
		Privacy:    userFromDB.Privacy,
		Membership: userFromDB.Membership,
		Preferences: &serverModels.UserPreferences{
			LanguageCode: userFromDB.Preferences.LanguageCode,
		},
		Friends: extractFriends(userFromDB),
	}

	return returnedUser
}

func extractFriends(user *models.User) []serverModels.Friend {
	var friends []serverModels.Friend
	for _, friend := range user.Friends {
		returnableFriend := serverModels.Friend{
			Username: friend.Username,
		}
		friends = append(friends, returnableFriend)
	}
	return friends
}

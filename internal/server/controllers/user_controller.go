package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/factories"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// CreateUser adds a new user to the database
func (env *Env) CreateUser(_ *gin.Context, user *serverModels.NewUser) error {
	userLogger := log.WithFields(log.Fields{
		"username": user.Username,
	})
	userLogger.Debug("Trying to add new user.")
	userService := biz.UserService{DB: env.DB}
	err := userService.Add(user.Username, user.Membership, user.Admin)

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

	userService := biz.UserService{DB: env.DB}
	usernames, err := userService.GetAll(adminFilter, privacyFilter, membershipFilter)

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

	userService := biz.UserService{DB: env.DB}
	userFromDB, err := userService.Get(params.Username)
	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("User does not exist."))
		return &serverModels.User{}, errors.NotFoundf("The user %s", params.Username)
	}

	returnedUser := newUser(userFromDB)
	return returnedUser, nil
}

// RemoveUser removes a user
func (env *Env) RemoveUser(_ *gin.Context, user *serverModels.UserParams) error {
	userLogger := log.WithFields(log.Fields{
		"username": user.Username,
	})
	userLogger.Debug("Trying to add new user.")

	userService := biz.UserService{DB: env.DB}
	err := userService.Remove(user.Username)
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
		Friends: extractUserFriends(userFromDB),
	}

	return returnedUser
}

func extractUserFriends(user *models.User) []serverModels.Friend {
	var friends []serverModels.Friend
	for _, friend := range user.Friends {
		returnableFriend := serverModels.Friend{
			Username: friend.Username,
		}
		friends = append(friends, returnableFriend)
	}
	return friends
}

// GetUserPools returns all the user's questions pool.
func (env *Env) GetUserPools(_ *gin.Context, params *serverModels.UserParams) ([]serverModels.QuestionPool, error) {
	userLogger := log.WithFields(log.Fields{
		"username": params.Username,
	})
	userLogger.Debug("Trying to get user pools.")
	userService := biz.UserService{DB: env.DB}
	pools, err := userService.GetPools(params.Username)

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("User does not exist."))
		return []serverModels.QuestionPool{}, errors.NotFoundf("The user %s", params.Username)
	}

	userPools, err := env.getUserPools(pools)
	return userPools, err
}

func (env *Env) getUserPools(questionPools []models.QuestionPool) ([]serverModels.QuestionPool, error) {
	var pools []serverModels.QuestionPool

	for _, pool := range questionPools {
		game, err := factories.GetGame(pool.GameName)
		if err != nil {
			env.Logger.Errorf("Unknown game type %s", pool.GameName)
		}

		questionPoolQuestions, err := game.NewQuestionPool(pool.Questions)
		if err != nil {
			return []serverModels.QuestionPool{}, err
		}
		newPool := serverModels.QuestionPool{
			PoolName:  pool.PoolName,
			GameName:  pool.GameName,
			Privacy:   pool.Privacy,
			Questions: questionPoolQuestions,
		}
		pools = append(pools, newPool)
	}

	return pools, nil
}

// GetUserStories returns all the user's stories.
func (env *Env) GetUserStories(_ *gin.Context, params *serverModels.UserParams) ([]serverModels.Story, error) {
	userLogger := log.WithFields(log.Fields{
		"username": params.Username,
	})
	userLogger.Debug("Trying to get user stories.")

	userService := biz.UserService{DB: env.DB}
	userStories, err := userService.GetUserStories(params.Username)
	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("User does not exist."))
		return []serverModels.Story{}, errors.NotFoundf("The user %s", params.Username)
	}

	stories := env.getUserStories(userStories)
	return stories, nil
}

func (env *Env) getUserStories(userStories []models.Story) []serverModels.Story {
	stories := []serverModels.Story{}

	for _, story := range userStories {
		game, err := factories.GetGame(story.GameName)
		if err != nil {
			env.Logger.Errorf("Unknown game type %s", story.GameName)
		}
		newStory := game.NewStory(story)
		stories = append(stories, newStory)
	}

	return stories
}

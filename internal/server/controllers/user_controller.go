package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"

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

// GetAllUserPools returns all the user's questions pool.
func (env *Env) GetAllUserPools(_ *gin.Context, params *serverModels.UserParams) ([]serverModels.QuestionPool, error) {
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

// GetUserPool returns a single question pool for a specified user.
func (env *Env) GetUserPool(
	_ *gin.Context,
	params *serverModels.ExistingQuestionPoolParams,
) (serverModels.QuestionPool, error) {
	userLogger := log.WithFields(log.Fields{
		"username":  params.Username,
		"pool_name": params.PoolName,
	})
	userLogger.Debug("Trying to get a single user pool.")
	userService := biz.UserService{DB: env.DB}
	userPool, err := userService.GetPool(params.Username, params.PoolName)

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Something went wrong, most likely because username or pool name does not exist.")
		return serverModels.QuestionPool{}, err
	}

	singlePool, err := env.newQuestionPool(userPool)
	return singlePool, err
}

func (env *Env) getUserPools(questionPools []models.QuestionPool) ([]serverModels.QuestionPool, error) {
	var pools []serverModels.QuestionPool

	for _, pool := range questionPools {
		newPool, err := env.newQuestionPool(pool)
		if err != nil {
			return []serverModels.QuestionPool{}, err
		}

		pools = append(pools, newPool)
	}

	return pools, nil
}

func (env *Env) newQuestionPool(pool models.QuestionPool) (serverModels.QuestionPool, error) {
	game, err := factories.GetGame(pool.GameName)
	if err != nil {
		env.Logger.Errorf("Unknown game %s", pool.GameName)
		return serverModels.QuestionPool{}, err
	}

	questionPoolQuestions, err := game.NewQuestionPool(pool.Questions)
	if err != nil {
		return serverModels.QuestionPool{}, err
	}

	newPool := serverModels.QuestionPool{
		PoolName:     pool.PoolName,
		GameName:     pool.GameName,
		LanguageCode: pool.LanguageCode,
		Privacy:      pool.Privacy,
		Questions:    questionPoolQuestions,
	}
	return newPool, nil
}

// AddUserPool adds a new user pool for an existing user.
func (env *Env) AddUserPool(
	_ *gin.Context,
	input *serverModels.QuestionPoolInput,
) (struct{}, error) {
	userLogger := log.WithFields(log.Fields{
		"username": input.UserParams.Username,
	})
	userLogger.Debug("Trying to add question pool.")
	userService := biz.UserService{DB: env.DB}

	languageCode := input.NewQuestionPool.LanguageCode
	if languageCode == "" {
		languageCode = "en"
	}

	_, err := language.Parse(languageCode)
	if err != nil {
		log.Errorf("failed to parse language code %s, err %s", languageCode, err)
		return struct{}{}, errors.BadRequestf("invalid language code %s", languageCode)
	}

	err = userService.AddPool(
		input.UserParams.Username,
		input.NewQuestionPool.PoolName,
		languageCode,
		input.NewQuestionPool.GameName,
		input.NewQuestionPool.Privacy,
	)

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Could not add question pool.")
	}

	return struct{}{}, err
}

// RemoveUserPool removes an existing question pool (for a specific user).
func (env *Env) RemoveUserPool(
	_ *gin.Context,
	input *serverModels.ExistingQuestionPoolParams,
) ([]serverModels.QuestionPool, error) {
	userLogger := log.WithFields(log.Fields{
		"username":  input.UserParams.Username,
		"pool_name": input.PoolParams.PoolName,
	})
	userLogger.Debug("Trying to add question pool.")
	userService := biz.UserService{DB: env.DB}

	err := userService.RemovePool(
		input.UserParams.Username,
		input.PoolParams.PoolName,
	)

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Could not remove question pool.")
	}

	return []serverModels.QuestionPool{}, err
}

// UpdateUserPool adds or removes a question from an existing question pool (for a specific user).
func (env *Env) UpdateUserPool(
	_ *gin.Context,
	input *serverModels.UpdateQuestionPoolInput,
) (struct{}, error) {
	userLogger := log.WithFields(log.Fields{
		"username":  input.UserParams.Username,
		"pool_name": input.PoolParams.PoolName,
		"operation": input.UpdateQuestionPool.Operation,
	})
	userLogger.Debug("Trying to update question pool.")

	questionToUpdate := env.newGenericQuestion(input.UpdateQuestionPool.NewQuestion)
	userService := biz.UserService{DB: env.DB}

	err := userService.UpdatePool(
		input.UserParams.Username,
		input.PoolParams.PoolName,
		input.UpdateQuestionPool.Operation,
		questionToUpdate,
	)

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Could update question pool.")
	}

	return struct{}{}, err
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

	stories, err := env.getUserStories(userStories)
	return stories, err
}

func (env *Env) getUserStories(userStories []models.Story) ([]serverModels.Story, error) {
	stories := []serverModels.Story{}

	for _, story := range userStories {
		game, err := factories.GetGame(story.GameName)
		if err != nil {
			env.Logger.Errorf("unknown game %s", story.GameName)
			return []serverModels.Story{}, err
		}
		newStory, err := game.NewStory(story)
		if err != nil {
			return []serverModels.Story{}, err
		}
		stories = append(stories, newStory)
	}

	return stories, nil
}

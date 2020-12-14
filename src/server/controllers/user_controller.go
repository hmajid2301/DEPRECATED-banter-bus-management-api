package controllers

import (
	"banter-bus-server/src/core"
	"banter-bus-server/src/core/models"
	serverModels "banter-bus-server/src/server/models"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

// CreateUser adds a new user to the database
func CreateUser(_ *gin.Context, user *serverModels.NewUser) error {
	userLogger := log.WithFields(log.Fields{
		"username": user.Username,
	})
	userLogger.Debug("Trying to add new user.")
	err := core.AddUser(user.Username, user.Membership, user.Admin)

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
func GetAllUsers(_ *gin.Context, params *serverModels.ListUserParams) ([]string, error) {
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

	usernames, err := core.GetAllUsers(adminFilter, privacyFilter, membershipFilter)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get user types.")
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
func GetUser(_ *gin.Context, params *serverModels.UserParams) (*serverModels.User, error) {
	userLogger := log.WithFields(log.Fields{
		"username": params.Username,
	})
	userLogger.Debug("Trying to add new user type.")

	userFromDB, err := core.GetUser(params.Username)

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("User doesn't exists"))
		return &serverModels.User{}, errors.NotFoundf("The user %s", params.Username)
	}

	returnedUser := newUser(userFromDB)

	return returnedUser, nil
}

// RemoveUser removes a user
func RemoveUser(_ *gin.Context, user *serverModels.UserParams) error {
	userLogger := log.WithFields(log.Fields{
		"username": user.Username,
	})
	userLogger.Debug("Trying to add new user type.")

	err := core.RemoveUser(user.Username)
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
		Stories: extractUserStories(userFromDB),
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

func extractUserStories(user *models.User) []serverModels.Story {
	var stories []serverModels.Story

	for _, story := range user.Stories {
		var answers []serverModels.StoryAnswer
		for _, storyAnswer := range story.Answers {
			returnableStoryAnswer := serverModels.StoryAnswer{
				Answer: storyAnswer.Answer,
				Votes:  storyAnswer.Votes,
			}
			answers = append(answers, returnableStoryAnswer)
		}

		returnableStory := serverModels.Story{
			GameName: story.GameName,
			Question: story.Question,
			Answers:  answers,
		}

		stories = append(stories, returnableStory)
	}

	return stories
}

// GetUserPools returns all the users questions pool.
func GetUserPools(_ *gin.Context, params *serverModels.UserParams) ([]serverModels.QuestionPool, error) {
	userLogger := log.WithFields(log.Fields{
		"username": params.Username,
	})
	userLogger.Debug("Trying to get user pools.")
	pools, err := core.GetUserPools(params.Username)

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("User doesn't exists"))
		return []serverModels.QuestionPool{}, errors.NotFoundf("The user %s", params.Username)
	}

	userPools := getUserPools(pools)
	return userPools, nil
}

func getUserPools(questionPools []models.QuestionPool) []serverModels.QuestionPool {
	var pools []serverModels.QuestionPool

	for _, pool := range questionPools {
		newQuestionsList := newQuestionPoolGenericQuestionList(pool.Questions.([]models.GenericQuestion))
		newPool := serverModels.QuestionPool{
			PoolName:  pool.PoolName,
			GameName:  pool.GameName,
			Privacy:   pool.Privacy,
			Questions: newQuestionsList,
		}
		pools = append(pools, newPool)
	}

	return pools
}

func newQuestionPoolGenericQuestionList(questions []models.GenericQuestion) []serverModels.GenericQuestion {
	var newQuestionsList []serverModels.GenericQuestion

	for _, question := range questions {
		newGenericQuestion := newQuestionPoolGenericQuestion(question)
		newQuestionsList = append(newQuestionsList, newGenericQuestion)
	}

	return newQuestionsList
}

func newQuestionPoolGenericQuestion(question models.GenericQuestion) serverModels.GenericQuestion {
	var group *serverModels.GenericQuestionGroup
	if question.Group != nil {
		group = &serverModels.GenericQuestionGroup{
			Name: question.Group.Name,
			Type: question.Group.Type,
		}
	}

	newQuestion := serverModels.GenericQuestion{
		Content: question.Content,
		Round:   question.Round,
		Group:   group,
	}

	return newQuestion
}

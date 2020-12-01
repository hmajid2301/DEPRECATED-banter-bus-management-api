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
	userLogger.Debug("Trying to add new user type.")
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

func newUser(userFromDB *models.User) *serverModels.User {
	returnedUser := &serverModels.User{
		Username:   userFromDB.Username,
		Admin:      userFromDB.Admin,
		Privacy:    userFromDB.Privacy,
		Membership: userFromDB.Membership,
		Preferences: &serverModels.UserPreferences{
			LanguageCode: userFromDB.Preferences.LanguageCode,
		},
		Friends:       extractUserFriends(userFromDB),
		Stories:       extractUserStories(userFromDB),
		QuestionPools: extractUserQuestionPools(userFromDB),
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

func extractUserQuestionPools(user *models.User) []serverModels.QuestionPool {
	var pools []serverModels.QuestionPool

	for range user.QuestionPools {
		// TODO implement properly when question pools are implemented properly
		pools = append(pools, serverModels.QuestionPool{})
	}
	return pools
}

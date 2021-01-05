package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
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
func (env *Env) GetUser(_ *gin.Context, params *serverModels.UserParams) (*serverModels.User, error) {
	userLogger := log.WithFields(log.Fields{
		"username": params.Username,
	})
	userLogger.Debug("Trying to add new user type.")

	userService := biz.UserService{DB: env.DB}
	userFromDB, err := userService.Get(params.Username)
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
func (env *Env) RemoveUser(_ *gin.Context, user *serverModels.UserParams) error {
	userLogger := log.WithFields(log.Fields{
		"username": user.Username,
	})
	userLogger.Debug("Trying to add new user type.")

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
		}).Warn(("User doesn't exists"))
		return []serverModels.QuestionPool{}, errors.NotFoundf("The user %s", params.Username)
	}

	userPools := env.getUserPools(pools)
	return userPools, nil
}

func (env *Env) getUserPools(questionPools []models.QuestionPool) []serverModels.QuestionPool {
	var pools []serverModels.QuestionPool

	for _, pool := range questionPools {
		questionPoolQuestions := env.newQuestionPoolQuestions(pool.GameName, pool.Questions)
		newPool := serverModels.QuestionPool{
			PoolName:  pool.PoolName,
			GameName:  pool.GameName,
			Privacy:   pool.Privacy,
			Questions: questionPoolQuestions,
		}

		pools = append(pools, newPool)
	}

	return pools
}

func (env *Env) newQuestionPoolQuestions(gameName string, questions interface{}) serverModels.QuestionPoolQuestions {
	var questionPool serverModels.QuestionPoolQuestions
	switch gameName {
	case "quibly":
		questionPool = newQuiblyQuestionPool(questions)
	case "fibbing_it":
		questionPool = newFibbingItQuestionPool(questions)
	case "drawlosseum":
		questionPool = newDrawlosseumQuestionPool(questions)
	default:
		env.Logger.Warnf("Invalid game %s", gameName)
	}

	return questionPool
}

func newQuiblyQuestionPool(questions interface{}) serverModels.QuestionPoolQuestions {
	quiblyQuestions := questions.(models.QuiblyQuestionsPool)

	var newQuiblyQuestionPool serverModels.QuiblyQuestionsPool
	newQuiblyQuestionPool.Pair = quiblyQuestions.Pair
	newQuiblyQuestionPool.Group = quiblyQuestions.Group
	newQuiblyQuestionPool.Answers = quiblyQuestions.Answers
	return serverModels.QuestionPoolQuestions{Quibly: newQuiblyQuestionPool}
}

func newFibbingItQuestionPool(questions interface{}) serverModels.QuestionPoolQuestions {
	fibbingItQuestions := questions.(models.FibbingItQuestionsPool)

	var newFibbingItQuestionPool serverModels.FibbingItQuestionsPool
	newFibbingItQuestionPool.Opinion = fibbingItQuestions.Opinion
	newFibbingItQuestionPool.Likely = fibbingItQuestions.Likely
	newFibbingItQuestionPool.FreeForm = fibbingItQuestions.FreeForm
	return serverModels.QuestionPoolQuestions{FibbingIt: newFibbingItQuestionPool}
}

func newDrawlosseumQuestionPool(questions interface{}) serverModels.QuestionPoolQuestions {
	drawlosseumQuestions := questions.(models.DrawlosseumQuestionsPool)

	var newDrawlosseumQuestionPool serverModels.DrawlosseumQuestionsPool
	newDrawlosseumQuestionPool.Drawings = drawlosseumQuestions.Drawings
	return serverModels.QuestionPoolQuestions{Drawlosseum: newDrawlosseumQuestionPool}
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
		}).Warn(("User doesn't exists"))
		return []serverModels.Story{}, errors.NotFoundf("The user %s", params.Username)
	}

	stories := getUserStories(userStories)
	return stories, nil
}

func getUserStories(userStories []models.Story) []serverModels.Story {
	stories := []serverModels.Story{}

	for _, story := range userStories {
		switch story.GameName {
		case "quibly":
			newStory := newQuiblyStory(story)
			stories = append(stories, newStory)
		case "fibbing_it":
			newStory := newFibbingItStory(story)
			stories = append(stories, newStory)
		case "drawlosseum":
			newStory := newDrawlosseumStory(story)
			stories = append(stories, newStory)
		}
	}

	return stories
}

func newQuiblyStory(story models.Story) serverModels.Story {
	quiblyAnswers := newAnswersQuibly(story.Answers)
	newQuiblyStory := serverModels.Story{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswers: serverModels.StoryAnswers{
			Quibly: quiblyAnswers,
		},
	}
	return newQuiblyStory
}

func newAnswersQuibly(answers interface{}) []serverModels.StoryQuibly {
	quiblyAnswers := answers.([]models.StoryQuibly)

	var newAnswersQuibly []serverModels.StoryQuibly
	for _, answer := range quiblyAnswers {
		newAnswer := serverModels.StoryQuibly{
			Nickname: answer.Nickname,
			Answer:   answer.Answer,
			Votes:    answer.Votes,
		}
		newAnswersQuibly = append(newAnswersQuibly, newAnswer)
	}

	return newAnswersQuibly
}

func newFibbingItStory(story models.Story) serverModels.Story {
	fibbingItAnswers := newAnswersFibbingIt(story.Answers)
	newFibbingItStory := serverModels.Story{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswers: serverModels.StoryAnswers{
			FibbingIt: fibbingItAnswers,
		},
	}
	return newFibbingItStory
}

func newAnswersFibbingIt(answers interface{}) []serverModels.StoryFibbingIt {
	fibbingItAnswers := answers.([]models.StoryFibbingIt)

	var newAnswersFibbingIt []serverModels.StoryFibbingIt
	for _, answer := range fibbingItAnswers {
		newAnswer := serverModels.StoryFibbingIt{
			Nickname: answer.Nickname,
			Answer:   answer.Answer,
		}
		newAnswersFibbingIt = append(newAnswersFibbingIt, newAnswer)
	}

	return newAnswersFibbingIt
}

func newDrawlosseumStory(story models.Story) serverModels.Story {
	drawlosseumAnswers := newAnswersDrawlosseum(story.Answers)
	newDrawlosseumStory := serverModels.Story{
		Question: story.Question,
		Nickname: story.Nickname,
		StoryAnswers: serverModels.StoryAnswers{
			Drawlosseum: drawlosseumAnswers,
		},
	}
	return newDrawlosseumStory
}

func newAnswersDrawlosseum(answers interface{}) []serverModels.StoryDrawlosseum {
	drawlosseumAnswers := answers.([]models.StoryDrawlosseum)

	var newAnswersDrawlosseum []serverModels.StoryDrawlosseum
	for _, answer := range drawlosseumAnswers {
		newAnswer := serverModels.StoryDrawlosseum{
			Color: answer.Color,
			Start: serverModels.DrawlosseumDrawingPoint{
				X: answer.Start.X,
				Y: answer.Start.Y,
			},
			End: serverModels.DrawlosseumDrawingPoint{
				X: answer.End.X,
				Y: answer.End.Y,
			},
		}
		newAnswersDrawlosseum = append(newAnswersDrawlosseum, newAnswer)
	}

	return newAnswersDrawlosseum
}

package service

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
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
		Username:      u.Username,
		Admin:         admin,
		Privacy:       "private",
		Membership:    membership,
		Preferences:   &models.UserPreferences{},
		Friends:       []models.Friend{},
		QuestionPools: []models.QuestionPool{},
	}

	inserted, err := user.Add(u.DB)
	if !inserted {
		return errors.Errorf("failed to add user %s", u.Username)
	}
	return err
}

// Remove removes a user from the database
func (u *UserService) Remove() error {
	exists, err := u.doesUserExist()
	if !exists {
		return errors.NotFoundf("the user %s", u.Username)
	} else if err != nil {
		return err
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

// GetPools gets a specific user's (all of them) question pools
func (u *UserService) GetPools() ([]models.QuestionPool, error) {
	user, err := u.Get()
	if err != nil {
		return []models.QuestionPool{}, err
	}
	pools := user.QuestionPools
	return pools, nil
}

// GetPool gets a single question pool from a user
func (u *UserService) GetPool(poolName string) (models.QuestionPool, error) {
	var (
		filter      = map[string]string{"username": u.Username}
		parentField = "question_pools"
		pool        models.QuestionPools
		condition   = []string{"$$this.pool_name", poolName}
	)

	err := pool.Get(u.DB, filter, parentField, condition)
	if (err != nil) || (len(pool) < 1) {
		return models.QuestionPool{}, errors.NotFoundf("user %s and pool %s", u.Username, poolName)
	}

	return pool[0], nil
}

// AddPool add question pool to a user
func (u *UserService) AddPool(
	poolName string,
	langCode string,
	gameName string,
	privacy string,
) error {
	exists, err := u.doesUserExist()
	if !exists {
		return errors.NotFoundf("username %s", u.Username)
	} else if err != nil {
		return err
	}

	currPool, _ := u.GetPool(poolName)
	if (currPool != models.QuestionPool{}) {
		return errors.AlreadyExistsf("pool %s", poolName)
	}

	game, err := games.GetGame(gameName)
	if err != nil {
		return err
	}

	emptyQuestions := game.NewQuestionPool()
	pool := models.NewPoolQuestion{
		"question_pools": {
			PoolName:     poolName,
			GameName:     gameName,
			LanguageCode: langCode,
			Privacy:      privacy,
			Questions:    emptyQuestions,
		},
	}

	filter := map[string]string{"username": u.Username}
	inserted, err := pool.AddToList(u.DB, filter)

	if !inserted && err == nil {
		return errors.Errorf("failed to add new pool %s", poolName)
	}
	return err
}

// RemovePool deletes a question pool from an user
func (u *UserService) RemovePool(poolName string) error {
	exists, _ := u.doesUserExist()
	if !exists {
		return errors.NotFoundf("username %s", u.Username)
	}

	_, err := u.GetPool(poolName)
	if err != nil {
		return err
	}

	pool := models.UpdateUserObject{
		"question_pools": map[string]string{
			"pool_name": poolName,
		},
	}

	filter := map[string]string{"username": u.Username}
	removed, err := pool.RemoveFromList(u.DB, filter)

	if !removed {
		return errors.Errorf("failed to remove new pool %s", poolName)
	}
	return err
}

// UpdatePool adds or removes questions from an existing question pool (for a specific user).
// It checks if the user exists, if the question pool exists and if the question already exists as well and handles
// those cases as required.
func (u *UserService) UpdatePool(poolName string, operation string, question models.GenericQuestion) error {
	pool, err := u.GetPool(poolName)
	if err != nil {
		return err
	}

	err = validateQuestion(pool.GameName, question)
	if err != nil {
		return err
	}

	genericQuestions, err := getGenericQuestions(pool.GameName, pool.Questions)
	if err != nil {
		return err
	}

	updatedQuestion, err := updatedPoolQuestion(pool.GameName, question)
	if err != nil {
		return err
	}

	filter := map[string]string{
		"username":                 u.Username,
		"question_pools.pool_name": poolName,
	}

	var updated bool
	exists := questionExist(genericQuestions, question)
	if operation == "add" {
		updated, err = u.addQuestion(exists, question.Content, filter, updatedQuestion)
	} else if operation == "remove" {
		updated, err = u.removeQuestion(exists, question.Content, filter, updatedQuestion)
	}

	if !updated && err == nil {
		return errors.Errorf("failed to update (%s) pool question %s", operation, poolName)
	}
	return err
}

func (u *UserService) addQuestion(
	exists bool,
	content string,
	filter map[string]string,
	question models.UpdateUserObject,
) (bool, error) {
	if exists {
		return false, errors.AlreadyExistsf("question '%s'", content)
	}
	updated, err := question.AddToList(u.DB, filter)
	return updated, err
}

func (u *UserService) removeQuestion(
	exists bool,
	content string,
	filter map[string]string,
	question models.UpdateUserObject,
) (bool, error) {
	if !exists {
		return false, errors.NotFoundf("question '%s'", content)
	}
	updated, err := question.RemoveFromList(u.DB, filter)
	return updated, err
}

func questionExist(currQuestions []models.GenericQuestion, question models.GenericQuestion) bool {
	if question.Group != nil {
		group := &models.GenericQuestionGroup{
			Name: "",
			Type: "",
		}

		if cmp.Equal(question.Group, group) {
			question.Group = nil
		}
	}

	for _, currQuestion := range currQuestions {
		currQuestion.LanguageCode = ""
		question.LanguageCode = ""

		if cmp.Equal(question, currQuestion) {
			return true
		}
	}

	return false
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

func getGenericQuestions(gameName string, questions models.QuestionPoolType) ([]models.GenericQuestion, error) {
	game, err := games.GetGame(gameName)
	if err != nil {
		return []models.GenericQuestion{}, err
	}

	genericQuestions, err := game.QuestionPoolToGenericQuestions(questions)
	return genericQuestions, err
}

func updatedPoolQuestion(gameName string, question models.GenericQuestion) (models.UpdateUserObject, error) {
	endPath, err := getQuestionPath(gameName, question)
	path := fmt.Sprintf("question_pools.$.%s", endPath)
	updatedQuestion := models.UpdateUserObject{
		path: question.Content,
	}

	return updatedQuestion, err
}

func getQuestionPath(gameName string, question models.GenericQuestion) (string, error) {
	game, err := games.GetGame(gameName)
	if err != nil {
		return "", err
	}

	path := game.GetQuestionPath(question)
	return path, nil
}

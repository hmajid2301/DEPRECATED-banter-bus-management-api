package biz

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
)

// UserService is struct data required by all user service functions.
type UserService struct {
	DB models.Repository
}

// Add is used to add a new user type
func (u *UserService) Add(username string, membership string, admin *bool) error {
	if u.doesUserExist(username) {
		return errors.AlreadyExistsf("The user %s", username)
	}

	var user = models.User{
		Username:      username,
		Admin:         admin,
		Privacy:       "public",
		Membership:    membership,
		Preferences:   &models.UserPreferences{},
		Friends:       []models.Friend{},
		Stories:       []models.Story{},
		QuestionPools: []models.QuestionPool{},
	}

	inserted, err := user.Add(u.DB)
	if !inserted {
		logger := log.WithFields(log.Fields{
			"username": username,
			"err":      err,
		})
		logger.Errorf("Failed to add new user %s", username)
		return errors.Errorf("failed to add new user %s", username)
	}
	return err
}

// GetAll is used to get usernames of all users
func (u *UserService) GetAll(adminFilter *bool, privacyFilter *string, membershipFilter *string) ([]string, error) {
	users := models.Users{}

	err := users.Get(u.DB)
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
		user   = &models.User{}
	)

	err := user.Get(u.DB, filter)
	if err != nil {
		return &models.User{}, errors.NotFoundf("user %s, %s", username, err)
	}

	return user, nil
}

// GetPools gets a specific user's (all of them) question pools
func (u *UserService) GetPools(username string) ([]models.QuestionPool, error) {
	user, err := u.Get(username)
	if err != nil {
		return []models.QuestionPool{}, err
	}
	pools := user.QuestionPools
	return pools, nil
}

// GetPool gets a single question pool from a user
func (u *UserService) GetPool(username string, poolName string) (models.QuestionPool, error) {
	var (
		filter       = map[string]string{"username": username}
		parentField  = "question_pools"
		questionPool models.QuestionPools
		condition    = []string{"$$this.pool_name", poolName}
	)

	err := questionPool.Get(u.DB, filter, parentField, condition)
	if (err != nil) || (len(questionPool) < 1) {
		return models.QuestionPool{}, errors.NotFoundf("user %s and pool %s", username, poolName)
	}

	return questionPool[0], nil
}

// AddPool add question pool to a user
func (u *UserService) AddPool(
	username string,
	poolName string,
	languageCode string,
	gameName string,
	privacy string,
) error {
	exists := u.doesUserExist(username)
	if !exists {
		return errors.NotFoundf("username %s", username)
	}

	pool, _ := u.GetPool(username, poolName)
	if (pool != models.QuestionPool{}) {
		return errors.AlreadyExistsf("pool %s", poolName)
	}

	game, err := games.GetGame(gameName)
	if err != nil {
		return err
	}

	questionPoolType := game.NewQuestionPool()
	newPool := models.NewQuestionPool{
		"question_pools": {
			PoolName:     poolName,
			GameName:     gameName,
			LanguageCode: languageCode,
			Privacy:      privacy,
			Questions:    questionPoolType,
		},
	}

	filter := map[string]string{"username": username}
	inserted, err := newPool.AddToList(u.DB, filter)

	if !inserted && err == nil {
		return errors.Errorf("failed to add new pool %s", poolName)
	}
	return err
}

// RemovePool deletes a question pool from an user
func (u *UserService) RemovePool(username string, poolName string) error {
	exists := u.doesUserExist(username)
	if !exists {
		return errors.NotFoundf("username %s", username)
	}

	_, err := u.GetPool(username, poolName)
	if err != nil {
		return err
	}

	poolToRemove := models.UpdateUserObject{
		"question_pools": map[string]string{
			"pool_name": poolName,
		},
	}

	filter := map[string]string{"username": username}
	removed, err := poolToRemove.RemoveFromList(u.DB, filter)

	if !removed {
		return errors.Errorf("failed to remove new pool %s", poolName)
	}
	return err
}

// UpdatePool adds or removes questions from an existing question pool (for a specific user).
// It checks if the user exists, if the question pool exists and if the question already exists as well and handles
// those cases as required.
func (u *UserService) UpdatePool(
	username string,
	poolName string,
	operation string,
	questionToUpdate models.GenericQuestion,
) error {
	poolToUpdate, err := u.GetPool(username, poolName)
	if err != nil {
		return err
	}

	var updated bool
	game, err := games.GetGame(poolToUpdate.GameName)
	if err != nil {
		return err
	}

	err = game.ValidateQuestionInput(questionToUpdate)
	if err != nil {
		return err
	}

	genericQuestions, err := game.QuestionPoolToGenericQuestions(poolToUpdate.Questions)
	if err != nil {
		return err
	}

	partialUpdateQuestionPath := game.GetQuestionPath(questionToUpdate)
	fullQuestionPath := fmt.Sprintf("question_pools.$.%s", partialUpdateQuestionPath)
	updateQuestion := models.UpdateUserObject{
		fullQuestionPath: questionToUpdate.Content,
	}

	filter := map[string]string{
		"username":                 username,
		"question_pools.pool_name": poolName,
	}

	exists := doesQuestionExist(genericQuestions, questionToUpdate)
	if operation == "add" {
		updated, err = u.addQuestionToPool(exists, questionToUpdate.Content, filter, updateQuestion)
	} else if operation == "remove" {
		updated, err = u.removeQuestionFromPool(exists, questionToUpdate.Content, filter, updateQuestion)
	}

	if !updated && err == nil {
		return errors.Errorf("failed to update (%s) pool question %s", operation, poolName)
	}
	return err
}

func (u *UserService) addQuestionToPool(
	exists bool,
	content string,
	filter map[string]string,
	newQuestion models.UpdateUserObject,
) (bool, error) {
	if exists {
		return false, errors.AlreadyExistsf("question '%s'", content)
	}
	updated, err := newQuestion.AddToList(u.DB, filter)
	return updated, err
}

func (u *UserService) removeQuestionFromPool(
	exists bool,
	content string,
	filter map[string]string,
	questionToRemove models.UpdateUserObject,
) (bool, error) {
	if !exists {
		return false, errors.NotFoundf("question '%s'", content)
	}
	updated, err := questionToRemove.RemoveFromList(u.DB, filter)
	return updated, err
}

func doesQuestionExist(
	existingQuestions []models.GenericQuestion,
	question models.GenericQuestion,
) bool {
	if question.Group != nil {
		group := &models.GenericQuestionGroup{
			Name: "",
			Type: "",
		}
		if cmp.Equal(question.Group, group) {
			question.Group = nil
		}
	}

	for _, existingQuesion := range existingQuestions {
		existingQuesion.LanguageCode = ""
		question.LanguageCode = ""

		if cmp.Equal(question, existingQuesion) {
			return true
		}
	}

	return false
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
